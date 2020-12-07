use crate::domain::handling::HandlingEventType as DomainHandlingEventType;
pub use pb::handling_service_client::HandlingServiceClient;
pub use pb::handling_service_server::{HandlingService, HandlingServiceServer};
pub use pb::{HandlingEventType, RegisterHandlingEventRequest};
use std::convert::{From, TryFrom};
use std::str::FromStr;

mod pb {
    tonic::include_proto!("handling"); // The string specified here must match the proto package name
}

#[derive(Debug)]
pub struct ParseError;

impl std::error::Error for ParseError {}

impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Parse error")
    }
}

impl FromStr for HandlingEventType {
    type Err = ParseError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "NotHandled" => Ok(Self::NotHandled),
            "Load" => Ok(Self::Load),
            "Unload" => Ok(Self::Unload),
            "Receive" => Ok(Self::Receive),
            "Claim" => Ok(Self::Claim),
            "Customs" => Ok(Self::Customs),
            _ => Err(Self::Err {}),
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
    type Error = &'static str;

    fn try_from(value: i32) -> Result<Self, Self::Error> {
        match value {
            0 => Ok(DomainHandlingEventType::NotHandled),
            1 => Ok(DomainHandlingEventType::Load),
            3 => Ok(DomainHandlingEventType::Unload),
            4 => Ok(DomainHandlingEventType::Receive),
            5 => Ok(DomainHandlingEventType::Claim),
            6 => Ok(DomainHandlingEventType::Customs),
            _ => Err("Can't convert to HandlingEventType"),
        }
    }
}
