use super::service::Service;
use crate::domain::handling::{HandlingEventType, TrackingID};
use crate::domain::location::UNLocode;
use crate::domain::voyage::VoyageNumber;
use crate::Error;
use chrono::prelude::*;
use log::debug;
use std::convert::{TryFrom, TryInto};
use std::time::SystemTime;
use tonic::{Code, Request, Response, Status};

use super::pb::{HandlingService, RegisterHandlingEventRequest};

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
        request: Request<RegisterHandlingEventRequest>,
    ) -> Result<Response<()>, Status> {
        let message = request.into_inner();
        debug!("{:?}", message);
        let completed = match message.completed {
            Some(prost_timestamp) => {
                let sys_time = SystemTime::try_from(prost_timestamp).unwrap();
                DateTime::<Utc>::from(sys_time)
            }
            None => Utc::now(),
        };

        let event_type_result: Result<HandlingEventType, Error> = message.event_type.try_into();
        let event_type = match event_type_result {
            Ok(etype) => etype,
            Err(err) => return Err(Status::new(Code::InvalidArgument, err.to_string())),
        };

        if let Err(error) = self
            .0
            .register_handling_event(
                completed,
                message.id as TrackingID,
                message.voyage_number as VoyageNumber,
                message.un_locode as UNLocode,
                event_type,
            )
            .await
        {
            return Err(Status::new(Code::Internal, error.to_string()));
        }
        Ok(Response::new(()))
    }
}
