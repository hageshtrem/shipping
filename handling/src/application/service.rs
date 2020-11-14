use crate::application::error::ErrInvalidArgument;
use crate::domain::handling::{
    HandlingEvent, HandlingEventFactory, HandlingEventType, TrackingID, VoyageNumber,
};
use crate::domain::location::UNLocode;
use crate::domain::{Repository, Result};
use chrono::prelude::*;

// pub trait Service {
//     fn register_handling_event(
//         &self,
//         id: TrackingID,
//         voyage_number: VoyageNumber,
//         un_locode: UNLocode,
//         event_type: HandlingEventType,
//     ) -> Result<(), Box<dyn std::error::Error>>;
// }

pub struct ServiceImpl<'a, R, F, H> {
    handling_event_repository: &'a R,
    handling_event_factory: &'a F,
    event_handler: &'a H,
}

impl<'a, R, F, H> ServiceImpl<'a, R, F, H>
where
    R: Repository<TrackingID, HandlingEvent> + 'a,
    F: HandlingEventFactory + 'a,
    H: EventHandler + 'a,
{
    pub fn new_service(
        handling_event_repository: &'a R,
        handling_event_factory: &'a F,
        event_handler: &'a H,
    ) -> Self {
        ServiceImpl {
            handling_event_repository,
            handling_event_factory,
            event_handler,
        }
    }
    pub fn register_handling_event(
        &self,
        completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<()> {
        match event_type {
            HandlingEventType::NotHandled
                if id.is_empty() || voyage_number.is_empty() || un_locode.is_empty() =>
            {
                return Err(Box::new(ErrInvalidArgument))
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
        self.event_handler.cargo_was_handled(e);
        Ok(())
    }
}

pub trait EventHandler {
    fn cargo_was_handled(&self, e: HandlingEvent);
}
