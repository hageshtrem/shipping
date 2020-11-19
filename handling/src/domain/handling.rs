#![allow(dead_code)]
use super::location::{Location, UNLocode};
use super::{Repository, Result};
use chrono::prelude::*;

#[derive(Debug, Clone)]
pub enum HandlingEventType {
    NotHandled,
    Load,
    Unload,
    Receive,
    Claim,
    Customs,
}

#[derive(Clone)]
pub struct HandlingActivity {
    r#type: HandlingEventType,
    location: UNLocode,
    voyage_number: VoyageNumber,
}

#[derive(Clone)]
pub struct HandlingEvent {
    tracking_id: TrackingID,
    activity: HandlingActivity,
}

pub trait HandlingEventFactory {
    fn create_handling_event(
        &self,
        registered: DateTime<Utc>,
        completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<HandlingEvent>;
}

pub struct HandlingEventFactoryImpl<C, V, L> {
    cargo_repository: C,
    voyage_repository: V,
    location_repository: L,
}

impl<C, V, L> HandlingEventFactoryImpl<C, V, L> {
    pub fn new(cargo_repository: C, voyage_repository: V, location_repository: L) -> Self {
        HandlingEventFactoryImpl {
            cargo_repository,
            voyage_repository,
            location_repository,
        }
    }
}

impl<C, V, L> HandlingEventFactory for HandlingEventFactoryImpl<C, V, L>
where
    C: Repository<TrackingID, Cargo>,
    V: Repository<VoyageNumber, Voyage>,
    L: Repository<UNLocode, Location>,
{
    fn create_handling_event(
        &self,
        _registered: DateTime<Utc>,
        _completed: DateTime<Utc>,
        id: TrackingID,
        voyage_number: VoyageNumber,
        un_locode: UNLocode,
        event_type: HandlingEventType,
    ) -> Result<HandlingEvent> {
        self.cargo_repository.find(id.clone())?;
        // When creating a Receive event, the voyage number is not known.
        if !voyage_number.is_empty() {
            self.voyage_repository.find(voyage_number.clone())?;
        }
        self.location_repository.find(un_locode.clone())?;

        Ok(HandlingEvent {
            tracking_id: id,
            activity: HandlingActivity {
                r#type: event_type,
                location: un_locode,
                voyage_number: voyage_number,
            },
        })
    }
}

pub type TrackingID = String;
pub type VoyageNumber = String;

#[derive(Clone)]
pub struct Cargo {}
#[derive(Clone)]
pub struct Voyage {}
