use crate::application::integration_events::event_handlers::EventHandler;
use crate::application::pb::TypeName;
use crate::domain::Result;
use async_trait::async_trait;
use bytes::Bytes;
use lapin::{
    options::*, types::FieldTable, Channel, Connection, ConnectionProperties, ExchangeKind,
};
use log::info;
use prost::Message;
use std::convert::From;
use tokio::stream::StreamExt;

const EXCHANGE_NAME: &'static str = "shipping";
const QUEUE_NAME: &'static str = "handling.queue";

#[async_trait]
pub trait SubscribeManager {
    async fn subscribe<E, EH>(&mut self, eh: EH)
    where
        E: Message + TypeName + Default,
        EH: EventHandler<Event = E> + Send + 'static;
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
        E: Message + TypeName + Default,
        EH: EventHandler<Event = E> + Send + 'static,
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

        let mut consumer = self
            .channel
            .basic_consume(
                QUEUE_NAME,
                "",
                BasicConsumeOptions::default(),
                FieldTable::default(),
            )
            .await
            .unwrap();

        tokio::spawn(async move {
            while let Some(delivery) = consumer.next().await {
                info!("received message");
                let (_, delivery) = delivery.expect("error in consumer");
                let e: E = Message::decode(Bytes::from(delivery.data)).unwrap();
                eh.handle(e);
            }
        });
    }
}
