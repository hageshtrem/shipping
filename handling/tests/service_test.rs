use chrono::prelude::*;
use handling::application::service::{EventHandlerImpl, Service, ServiceImpl};
use handling::domain::handling::{Cargo, HandlingEventFactoryImpl, HandlingEventType, Voyage};
use handling::domain::{location, Repository};
use handling::infrastructure::inmem_repository::InmemRepository;

#[test]
fn service() {
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
    voyages.store("v001".to_string(), &Voyage {}).unwrap();
    let locations = InmemRepository::new();
    location::store_sample_locations(&locations).unwrap();
    let handling_events = InmemRepository::new();
    let unimp_event_handler = EventHandlerImpl {};
    let event_factory = HandlingEventFactoryImpl::new(cargos, voyages, locations);

    // create service instance
    let srv = ServiceImpl::new_service(handling_events, event_factory, unimp_event_handler);
    let res = srv.register_handling_event(
        Utc::now(),
        "001".to_string(),
        "v001".to_string(),
        "SESTO".to_string(),
        HandlingEventType::Load,
    );

    println!("{:?}", res);
    assert!(res.is_ok());
}
