use chrono::prelude::*;
use handling::application::service::{EventHandlerImpl, ServiceImpl};
use handling::domain::handling::{Cargo, HandlingEventFactoryImpl, HandlingEventType, Voyage};
use handling::domain::{location, Repository};
use handling::infrastructure::inmem_repository::InmemRepository;

#[test]
fn service() {
    // prerare dependencies
    let cargos = InmemRepository::new();
    let voyages = InmemRepository::new();
    let locations = InmemRepository::new();
    let handling_events = InmemRepository::new();
    let unimp_event_handler = EventHandlerImpl {};
    let event_factory = HandlingEventFactoryImpl::new(&cargos, &voyages, &locations);

    // store test data
    location::store_sample_locations(&locations).unwrap();
    cargos.store("001".to_string(), &Cargo {}).unwrap();
    voyages.store("v001".to_string(), &Voyage {}).unwrap();

    // create service instance
    let srv = ServiceImpl::new_service(&handling_events, &event_factory, &unimp_event_handler);
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
