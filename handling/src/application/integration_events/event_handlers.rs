use crate::application::pb::NewCargoBooked;
use crate::domain::{handling::Cargo, handling::TrackingID, Repository};
use std::convert::TryInto;

pub trait EventHandler {
    type Event;
    fn handle(&self, e: Self::Event);
}

pub struct NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    cargo_repository: T,
}

impl<T> NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    pub fn new(cargo_repository: T) -> Self {
        NewCargoBookedEventHandler { cargo_repository }
    }
}

impl<T> EventHandler for NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    type Event = NewCargoBooked;

    fn handle(&self, e: Self::Event) {
        let cargo: Cargo = e.try_into().unwrap();
        self.cargo_repository
            .store(cargo.tracking_id.clone(), &cargo)
            .unwrap();
    }
}
