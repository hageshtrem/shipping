use crate::application::integration_events::event_handlers::EventHandler;
use crate::application::pb::TypeName;
use crate::domain::Result;
use async_trait::async_trait;
use bytes::Bytes;
use lapin::{
    message::DeliveryResult, options::*, types::FieldTable, Channel, Connection,
    ConnectionProperties, ConsumerDelegate, ExchangeKind,
};
use log::info;
use prost::Message;
use std::convert::From;
use std::future::Future;
use std::pin::Pin;

const EXCHANGE_NAME: &'static str = "shipping";
const QUEUE_NAME: &'static str = "handling.queue";

#[async_trait]
pub trait SubscribeManager {
    async fn subscribe<E, EH>(&mut self, eh: EH)
    where
        E: Message + TypeName + Default + 'static,
        EH: EventHandler<Event = E> + Send + Sync + 'static;
}

pub struct EventBus {
    channel: Channel,
}

impl EventBus {
    pub async fn new(url: &str) -> Result<Self> {
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
impl SubscribeManager for EventBus {
    async fn subscribe<E, EH>(&mut self, eh: EH)
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
            .wait()
            .unwrap();

        let consumer = self
            .channel
            .basic_consume(
                QUEUE_NAME,
                "",
                BasicConsumeOptions::default(),
                FieldTable::default(),
            )
            .await
            .unwrap();
        consumer.set_delegate(Delegate(eh)).unwrap();
    }
}

struct Delegate<E, EH>(EH)
where
    E: Message + TypeName + Default + 'static,
    EH: EventHandler<Event = E> + Send + Sync + 'static;

impl<E, EH> ConsumerDelegate for Delegate<E, EH>
where
    Self: 'static,
    E: Message + TypeName + Default + 'static,
    EH: EventHandler<Event = E> + Send + Sync + 'static,
{
    fn on_new_delivery(
        &self,
        delivery: DeliveryResult,
    ) -> Pin<Box<dyn Future<Output = ()> + Send>> {
        let eh = self.0.clone();
        Box::pin(async move {
            info!("received message");
            let delivery = delivery.expect("error in consumer");
            if let Some((_, delivery)) = delivery {
                let e: E = Message::decode(Bytes::from(delivery.data.clone())).unwrap();
                eh.handle(e);
                delivery.ack(BasicAckOptions::default()).await.expect("ack");
            }
        })
    }
}
