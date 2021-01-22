use super::integration_events::EventService;
use crate::domain::handling::{HandlingEvent, HandlingEventFactory, HandlingEventType, TrackingID};
use crate::domain::{location::UNLocode, voyage::VoyageNumber, Repository};
use crate::Error;
use async_trait::async_trait;
use chrono::prelude::*;

#[async_trait]
pub trait Service {
    async fn register_handling_event(
        &self,
        completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<(), Error>;
}

pub struct ServiceImpl<R, F, H> {
    handling_event_repository: R,
    handling_event_factory: F,
    event_handler: H,
}

impl<R, F, H> ServiceImpl<R, F, H>
where
    R: Repository<TrackingID, HandlingEvent>,
    F: HandlingEventFactory,
    H: EventService,
{
    pub fn new_service(
        handling_event_repository: R,
        handling_event_factory: F,
        event_handler: H,
    ) -> Self {
        ServiceImpl {
            handling_event_repository,
            handling_event_factory,
            event_handler,
        }
    }
}

#[async_trait]
impl<R, F, H> Service for ServiceImpl<R, F, H>
where
    R: Repository<TrackingID, HandlingEvent>,
    F: HandlingEventFactory,
    H: EventService,
{
    async fn register_handling_event(
        &self,
        completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<(), Error> {
        match event_type {
            HandlingEventType::NotHandled
                if id.is_empty() || voyage_number.is_empty() || un_locode.is_empty() =>
            {
                return Err(Error::InvalidArgument)
            }
            _ => (),
        }

        let e = self.handling_event_factory.create_handling_event(
            Utc::now(),
            completed,
            id.clone(),
            voyage_number,
            un_locode,
            event_type,
        )?;

        self.handling_event_repository.store(id, &e)?;
        self.event_handler.cargo_was_handled(e).await?;
        Ok(())
    }
}
