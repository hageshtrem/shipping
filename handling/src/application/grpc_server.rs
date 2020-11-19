use super::service::Service;
use crate::domain::handling::{HandlingEventType, TrackingID, VoyageNumber};
use crate::domain::location::UNLocode;
use chrono::prelude::*;
use std::convert::TryFrom;
use std::time::SystemTime;
use tonic::{Code, Request, Response, Status};

pub use pb::handling_service_server::{HandlingService, HandlingServiceServer};

pub mod pb {
    tonic::include_proto!("handling"); // The string specified here must match the proto package name
}

#[derive(Debug, Default)]
pub struct HandlingServiceImpl<S: Service>(S);

impl<S: Service> HandlingServiceImpl<S> {
    pub fn new(service: S) -> Self {
        HandlingServiceImpl(service)
    }
}

#[tonic::async_trait]
impl<S: Service + Sync + Send + 'static> HandlingService for HandlingServiceImpl<S> {
    async fn register_handling_event(
        &self,
        request: Request<pb::RegisterHandlingEventRequest>,
    ) -> Result<Response<()>, Status> {
        let message = request.into_inner();
        let completed = match message.completed {
            Some(prost_timestamp) => {
                let sys_time = SystemTime::try_from(prost_timestamp).unwrap();
                DateTime::<Utc>::from(sys_time)
            }
            None => Utc::now(),
        };
        let event_type = match message.event_type {
            0 => HandlingEventType::NotHandled,
            1 => HandlingEventType::Load,
            2 => HandlingEventType::Unload,
            3 => HandlingEventType::Receive,
            4 => HandlingEventType::Claim,
            5 => HandlingEventType::Customs,
            _ => return Err(Status::new(Code::InvalidArgument, "event type is invalid")),
        } as HandlingEventType;

        if let Err(error) = self.0.register_handling_event(
            completed,
            message.id as TrackingID,
            message.voyage_number as VoyageNumber,
            message.un_locode as UNLocode,
            event_type,
        ) {
            return Err(Status::new(Code::Internal, error.to_string()));
        }
        Ok(Response::new(()))
    }
}
