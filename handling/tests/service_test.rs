use async_trait::async_trait;
use chrono::prelude::*;
use handling::application::integration_events::EventService;
use handling::application::service::{Service, ServiceImpl};
use handling::domain::handling::{
    Cargo, HandlingEvent, HandlingEventFactoryImpl, HandlingEventType,
};
use handling::domain::{location, voyage, Repository};
use handling::infrastructure::inmem_repository::InmemRepository;
use handling::Error;

#[test]
fn service() {
    tokio_test::block_on(async {
        // prerare dependencies
        let cargos = InmemRepository::new();
        cargos
            .store(
                "001".to_string(),
                &Cargo {
                    tracking_id: "001".to_string(),
                    origin: "AUMEL".to_string(),
                    destination: "SESTO".to_string(),
                    arrival_deadline: Utc::now(),
                },
            )
            .unwrap();
        let voyages = InmemRepository::new();
        voyage::populate_repository(&voyages).unwrap();
        let locations = InmemRepository::new();
        location::store_sample_locations(&locations).unwrap();
        let handling_events = InmemRepository::new();
        let event_factory = HandlingEventFactoryImpl::new(cargos, voyages, locations);

        // create service instance
        let srv = ServiceImpl::new_service(handling_events, event_factory, MocEventService {});
        let res = srv
            .register_handling_event(
                Utc::now(),
                "001".to_string(),
                "0100S".to_string(),
                "SESTO".to_string(),
                HandlingEventType::Load,
            )
            .await;

        println!("{:?}", res);
        assert!(res.is_ok());
    });
}

struct MocEventService;

#[async_trait]
impl EventService for MocEventService {
    async fn cargo_was_handled(&self, _: HandlingEvent) -> Result<(), Error> {
        Ok(())
    }
}
