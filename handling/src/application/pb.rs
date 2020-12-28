use crate::domain::handling::Cargo;
use crate::domain::handling::HandlingActivity as DomainHandlingActivity;
use crate::domain::handling::HandlingEvent as DomainHandlingEvent;
use crate::domain::handling::HandlingEventType as DomainHandlingEventType;
use crate::Error;
use chrono::prelude::*;
pub use pb::booking::{CargoDestinationChanged, NewCargoBooked};
pub use pb::handling::handling_service_client::HandlingServiceClient;
pub use pb::handling::handling_service_server::{HandlingService, HandlingServiceServer};
pub use pb::handling::{Activity, HandlingEvent, HandlingEventType, RegisterHandlingEventRequest};
use std::convert::{From, Into, TryFrom};
use std::str::FromStr;
use std::time::SystemTime;

mod pb {
    pub mod handling {
        tonic::include_proto!("handling"); // The string specified here must match the proto package name
    }
    pub mod booking {
        tonic::include_proto!("booking");
    }
    pub mod itinerary {
        tonic::include_proto!("itinerary");
    }
}

impl FromStr for HandlingEventType {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "NotHandled" => Ok(Self::NotHandled),
            "Load" => Ok(Self::Load),
            "Unload" => Ok(Self::Unload),
            "Receive" => Ok(Self::Receive),
            "Claim" => Ok(Self::Claim),
            "Customs" => Ok(Self::Customs),
            _ => Err(Self::Err::ParsingError),
        }
    }
}

impl From<HandlingEventType> for DomainHandlingEventType {
    fn from(value: HandlingEventType) -> Self {
        match value {
            HandlingEventType::NotHandled => DomainHandlingEventType::NotHandled,
            HandlingEventType::Load => DomainHandlingEventType::Load,
            HandlingEventType::Unload => DomainHandlingEventType::Unload,
            HandlingEventType::Receive => DomainHandlingEventType::Receive,
            HandlingEventType::Claim => DomainHandlingEventType::Claim,
            HandlingEventType::Customs => DomainHandlingEventType::Customs,
        }
    }
}

impl TryFrom<i32> for DomainHandlingEventType {
    type Error = Error;
    fn try_from(value: i32) -> Result<Self, Self::Error> {
        match value {
            0 => Ok(DomainHandlingEventType::NotHandled),
            1 => Ok(DomainHandlingEventType::Load),
            2 => Ok(DomainHandlingEventType::Unload),
            3 => Ok(DomainHandlingEventType::Receive),
            4 => Ok(DomainHandlingEventType::Claim),
            5 => Ok(DomainHandlingEventType::Customs),
            _ => Err(Error::ParsingError),
        }
    }
}

impl From<DomainHandlingEventType> for i32 {
    fn from(value: DomainHandlingEventType) -> Self {
        match value {
            DomainHandlingEventType::NotHandled => 0,
            DomainHandlingEventType::Load => 1,
            DomainHandlingEventType::Unload => 2,
            DomainHandlingEventType::Receive => 3,
            DomainHandlingEventType::Claim => 4,
            DomainHandlingEventType::Customs => 5,
        }
    }
}

impl TryFrom<NewCargoBooked> for Cargo {
    type Error = Error;
    fn try_from(value: NewCargoBooked) -> Result<Self, Self::Error> {
        let arrival_deadline = match value.arrival_deadline {
            Some(prost_timestamp) => {
                let sys_time = SystemTime::try_from(prost_timestamp).unwrap();
                DateTime::<Utc>::from(sys_time)
            }
            None => Utc::now(), // TODO
        };
        Ok(Cargo {
            tracking_id: value.tracking_id,
            origin: value.origin,
            destination: value.destination,
            arrival_deadline: arrival_deadline,
        })
    }
}

impl From<DomainHandlingActivity> for Activity {
    fn from(value: DomainHandlingActivity) -> Self {
        Activity {
            r#type: value.r#type.into(),
            location: value.location,
            voyage_number: value.voyage_number,
        }
    }
}

impl From<DomainHandlingEvent> for HandlingEvent {
    fn from(value: DomainHandlingEvent) -> Self {
        HandlingEvent {
            tracking_id: value.tracking_id,
            activity: Some(value.activity.into()),
        }
    }
}

pub trait TypeName {
    fn name() -> &'static str;
}

impl TypeName for NewCargoBooked {
    fn name() -> &'static str {
        "NewCargoBooked"
    }
}

impl TypeName for CargoDestinationChanged {
    fn name() -> &'static str {
        "CargoDestinationChanged"
    }
}

impl TypeName for HandlingEvent {
    fn name() -> &'static str {
        "HandlingEvent"
    }
}
