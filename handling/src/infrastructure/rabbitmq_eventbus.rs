use crate::application::integration_events::{EventHandler, EventService};
use crate::application::pb::{HandlingEvent as PbHandlingEvent, TypeName};
use crate::domain::handling::HandlingEvent;
use crate::Error;
use async_trait::async_trait;
use bytes::Bytes;
use lapin::{
    message::DeliveryResult, options::*, types::FieldTable, BasicProperties, Channel, Connection,
    ConnectionProperties, ConsumerDelegate, ExchangeKind,
};
use log::{error, info};
use prost::Message;
use std::convert::From;
use std::future::Future;
use std::pin::Pin;

const EXCHANGE_NAME: &'static str = "shipping";
const QUEUE_NAME: &'static str = "handling.queue";

type DynResult<T> = Result<T, Box<dyn std::error::Error>>;

#[async_trait]
pub trait SubscribeManager {
    async fn subscribe<E, EH>(&mut self, eh: EH) -> DynResult<()>
    where
        E: Message + TypeName + Default + 'static,
        EH: EventHandler<Event = E> + Send + Sync + 'static;
}

pub struct EventBus {
    channel: Channel,
}

impl EventBus {
    pub async fn new(url: &str) -> DynResult<Self> {
        let conn = Connection::connect(url, ConnectionProperties::default()).await?;
        info!("Connected");
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

        Ok(EventBus { channel })
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
                BasicProperties::default(),
            )
            .await?;
        Ok(())
    }
}

#[async_trait]
impl SubscribeManager for EventBus {
    async fn subscribe<E, EH>(&mut self, eh: EH) -> DynResult<()>
    where
        E: Message + TypeName + Default + 'static,
        EH: EventHandler<Event = E> + Send + Sync + 'static,
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
        let consumer = self
            .channel
            .basic_consume(
                QUEUE_NAME,
                "",
                BasicConsumeOptions::default(),
                FieldTable::default(),
            )
            .await?;
        consumer.set_delegate(Delegate(eh))?;
        Ok(())
    }
}

struct Delegate<E, EH>(EH)
where
    E: Message + TypeName + Default + 'static,
    EH: EventHandler<Event = E> + Send + Sync + 'static;

impl<E, EH> ConsumerDelegate for Delegate<E, EH>
where
    E: Message + TypeName + Default + 'static,
    EH: EventHandler<Event = E> + Send + Sync + 'static,
{
    fn on_new_delivery(
        &self,
        delivery: DeliveryResult,
    ) -> Pin<Box<dyn Future<Output = ()> + Send>> {
        let eh = self.0.clone();
        Box::pin(async move {
            match delivery {
                Ok(delivery) => {
                    if let Some((_, delivery)) = delivery {
                        let e: std::result::Result<E, _> =
                            Message::decode(Bytes::from(delivery.data.clone()));
                        match e {
                            Ok(e) => match eh.handle(e) {
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
                            Err(err) => error!("error while decoding message: {}", err),
                        };
                    }
                }
                Err(err) => error!("error while receiving message: {}", err),
            };
        })
    }
}
