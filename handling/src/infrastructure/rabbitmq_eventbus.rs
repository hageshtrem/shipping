use crate::application::integration_events::{EventHandler, EventService};
use crate::application::pb::{HandlingEvent as PbHandlingEvent, TypeName};
use crate::domain::handling::HandlingEvent;
use crate::Error;
use async_trait::async_trait;
use bytes::Bytes;
use futures_util::stream::StreamExt;
use lapin::{
    options::*, types::FieldTable, BasicProperties, Channel, Connection, ConnectionProperties,
    Consumer, ExchangeKind,
};
use log::{error, info};
use prost::Message;
use std::collections::HashMap;
use std::convert::From;
use std::future::Future;
use std::ops::Deref;
use std::pin::Pin;
use std::sync::Arc;
use tokio::sync::Mutex;

const EXCHANGE_NAME: &'static str = "shipping";
const QUEUE_NAME: &'static str = "handling.queue";

type DynResult<T> = Result<T, Box<dyn std::error::Error>>;

#[async_trait]
pub trait SubscribeManager {
    async fn subscribe<E, EH>(&mut self, eh: EH) -> DynResult<()>
    where
        E: Message + TypeName + Default + 'static,
        EH: EventHandler<E> + Send + Sync + 'static;
}

pub struct EventBus {
    channel: Channel,
    crt: ConsumerRT,
}

impl EventBus {
    pub async fn new(url: &str) -> DynResult<Self> {
        let conn = Connection::connect(url, ConnectionProperties::default()).await?;
        let channel = conn.create_channel().await?;
        channel
            .exchange_declare(
                EXCHANGE_NAME,
                ExchangeKind::Direct,
                ExchangeDeclareOptions {
                    durable: true,
                    ..ExchangeDeclareOptions::default()
                },
                FieldTable::default(),
            )
            .await?;
        channel
            .queue_declare(
                QUEUE_NAME,
                QueueDeclareOptions::default(),
                FieldTable::default(),
            )
            .await?;
        let consumer = channel
            .basic_consume(
                QUEUE_NAME,
                "",
                BasicConsumeOptions::default(),
                FieldTable::default(),
            )
            .await?;
        let crt = ConsumerRT::new();
        crt.process(consumer).await;

        Ok(EventBus { channel, crt })
    }
}

#[async_trait]
impl EventService for EventBus {
    async fn cargo_was_handled(&self, e: HandlingEvent) -> Result<(), Error> {
        info!("{:?} cargo {}", e.activity.r#type, e.tracking_id);
        let pb_event: PbHandlingEvent = e.into();
        let mut buf = vec![];
        pb_event.encode(&mut buf)?;
        let _confitm = self
            .channel
            .basic_publish(
                EXCHANGE_NAME,
                PbHandlingEvent::name(),
                BasicPublishOptions::default(),
                buf,
                BasicProperties::default().with_kind(PbHandlingEvent::name().into()),
            )
            .await?;
        Ok(())
    }
}

type HandlerFunc =
    Box<dyn Fn(Vec<u8>) -> Pin<Box<dyn Future<Output = Result<(), Error>> + Send>> + Send + Sync>;

struct ConsumerRT {
    handlers: Arc<Mutex<HashMap<String, HandlerFunc>>>,
}

impl ConsumerRT {
    fn new() -> Self {
        let handlers = Arc::new(Mutex::new(HashMap::new()));
        ConsumerRT { handlers }
    }

    async fn add_handler_func(&mut self, msg_type: String, func: HandlerFunc) {
        let mut data = self.handlers.lock().await;
        data.insert(msg_type, func);
    }

    async fn process(&self, mut consumer: Consumer) {
        let handlers = self.handlers.clone();
        tokio::spawn(async move {
            while let Some(delivery) = consumer.next().await {
                match delivery {
                    Ok((_, delivery)) => {
                        if let Some(dtype) = delivery.properties.kind() {
                            let data = handlers.lock().await;
                            match data.deref().get(dtype.as_str()) {
                                Some(func) => match func(delivery.data.clone()).await {
                                    Ok(_) => delivery
                                        .ack(BasicAckOptions::default())
                                        .await
                                        .expect("RabbitMQ ack error"),
                                    Err(err) => {
                                        error!("error while handling event: {}", err);
                                        delivery
                                            .acker
                                            .nack(BasicNackOptions::default())
                                            .await
                                            .expect("RabbitMQ nack error");
                                    }
                                },
                                None => error!("No registered handler for: {}", dtype),
                            }
                        } else {
                            error!("Could not get message type {}", delivery.delivery_tag);
                        }
                    }
                    Err(err) => error!("error while receiving message: {}", err),
                };
            }
            info!("End of message stream")
        });
    }
}

#[async_trait]
impl SubscribeManager for EventBus {
    async fn subscribe<E, EH>(&mut self, eh: EH) -> DynResult<()>
    where
        E: Message + TypeName + Default + 'static,
        EH: EventHandler<E> + Send + Sync + 'static,
    {
        self.channel
            .queue_bind(
                QUEUE_NAME,
                EXCHANGE_NAME,
                E::name(),
                QueueBindOptions::default(),
                FieldTable::default(),
            )
            .await?;
        let eh = Arc::new(Mutex::new(eh));
        self.crt
            .add_handler_func(
                E::name().to_string(),
                Box::new(move |data: Vec<u8>| {
                    let eh = eh.clone();
                    Box::pin(async move {
                        let e: std::result::Result<E, _> = Message::decode(Bytes::from(data));
                        match e {
                            Ok(e) => {
                                let data = eh.lock().await;
                                match data.deref().handle(e) {
                                    Ok(_) => Ok(()),
                                    Err(err) => {
                                        error!("error while handling event: {}", err);
                                        Err(err)
                                    }
                                }
                            }
                            Err(err) => {
                                error!("error while decoding message: {}", err);
                                Err(err.into())
                            }
                        }
                    })
                }),
            )
            .await;
        Ok(())
    }
}
