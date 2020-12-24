use crate::application::pb::{CargoDestinationChanged, NewCargoBooked};
use crate::domain::{handling::Cargo, handling::TrackingID, Repository};
use crate::Error;
use log::info;
use std::convert::TryInto;

pub trait EventHandler: Clone + Send {
    type Event;
    fn handle(&self, e: Self::Event) -> Result<(), Error>;
}

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

impl<T> EventHandler for NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    type Event = NewCargoBooked;
    fn handle(&self, e: Self::Event) -> Result<(), Error> {
        let cargo: Cargo = e.try_into()?;
        info!("New cargo booked {}", cargo.tracking_id);
        self.cargos.store(cargo.tracking_id.clone(), &cargo)?;
        Ok(())
    }
}

impl<T> Clone for NewCargoBookedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    fn clone(&self) -> Self {
        NewCargoBookedEventHandler {
            cargos: self.cargos.clone(),
        }
    }
}

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

impl<T> EventHandler for CargoDestinationChangedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    type Event = CargoDestinationChanged;
    fn handle(&self, e: Self::Event) -> Result<(), Error> {
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

impl<T> Clone for CargoDestinationChangedEventHandler<T>
where
    T: Repository<TrackingID, Cargo>,
{
    fn clone(&self) -> Self {
        CargoDestinationChangedEventHandler {
            cargos: self.cargos.clone(),
        }
    }
}
