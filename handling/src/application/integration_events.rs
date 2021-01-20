use crate::application::pb::{CargoDestinationChanged, NewCargoBooked};
use crate::domain::handling::{Cargo, HandlingEvent, TrackingID};
use crate::domain::Repository;
use crate::Error;
use async_trait::async_trait;
use log::info;
use std::convert::TryInto;

#[async_trait]
pub trait EventService: Send + Sync {
    async fn cargo_was_handled(&self, e: HandlingEvent) -> Result<(), Error>;
}

pub trait EventHandler<Event>: Clone + Send {
    fn handle(&self, e: Event) -> Result<(), Error>;
}

#[derive(Clone)]
pub struct NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    cargos: T,
}

impl<T> NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    pub fn new(cargos: T) -> Self {
        NewCargoBookedEventHandler { cargos }
    }
}

impl<T> EventHandler<NewCargoBooked> for NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    fn handle(&self, e: NewCargoBooked) -> Result<(), Error> {
        let cargo: Cargo = e.try_into()?;
        info!("New cargo booked {}", cargo.tracking_id);
        self.cargos.store(cargo.tracking_id.clone(), &cargo)?;
        Ok(())
    }
}

#[derive(Clone)]
pub struct CargoDestinationChangedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    cargos: T,
}

impl<T> CargoDestinationChangedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    pub fn new(cargos: T) -> Self {
        CargoDestinationChangedEventHandler { cargos }
    }
}

impl<T> EventHandler<CargoDestinationChanged> for CargoDestinationChangedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    fn handle(&self, e: CargoDestinationChanged) -> Result<(), Error> {
        info!(
            "Cargo {} destination changed {}",
            e.tracking_id, e.destination
        );
        let mut cargo = self.cargos.find(e.tracking_id)?;
        cargo.destination = e.destination;
        self.cargos.store(cargo.tracking_id.clone(), &cargo)?;
        Ok(())
    }
}
