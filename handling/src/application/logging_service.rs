use super::service::Service;
use crate::domain::handling::{HandlingEventType, TrackingID};
use crate::domain::{location::UNLocode, voyage::VoyageNumber};
use crate::Error;
use async_trait::async_trait;
use chrono::prelude::*;
use log::info;

pub struct LoggingService {
    next: Box<dyn Service + Send + Sync>,
}

impl LoggingService {
    pub fn new(next: Box<dyn Service + Send + Sync>) -> impl Service {
        Self { next }
    }
}

#[async_trait]
impl Service for LoggingService {
    async fn register_handling_event(
        &self,
        completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<(), Error> {
        let begin = Utc::now();
        let res = self
            .next
            .register_handling_event(completed, id.clone(), voyage_number, un_locode, event_type)
            .await;
        let err = match &res {
            Ok(_) => "".to_string(),
            Err(err) => format!("{:?}", err),
        };
        info!(
            "method: register_handling_event, id: {}, err: {}, took: {}",
            id,
            err,
            Utc::now().signed_duration_since(begin)
        );
        res
    }
}
