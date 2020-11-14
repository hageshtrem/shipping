use crate::domain::location::{Location, UNLocode};
use crate::domain::{Repository, Result};
use std::collections::HashMap;
use std::ops::{Deref, DerefMut};
use std::sync::{Arc, Mutex};
use std::{error::Error, fmt};

struct InmemLocationRepository(Arc<Mutex<HashMap<UNLocode, Location>>>);

impl InmemLocationRepository {
    fn new() -> Self {
        InmemLocationRepository(Arc::new(Mutex::new(HashMap::new())))
    }
}

impl Repository<UNLocode, Location> for InmemLocationRepository {
    fn store(&self, un_locode: UNLocode, location: &Location) -> Result<()> {
        let r = self.0.clone();
        let mut data = r.lock().unwrap();
        data.deref_mut().insert(un_locode, location.clone());
        Ok(())
    }

    fn find(&self, un_locode: UNLocode) -> Result<Location> {
        let r = self.0.clone();
        let data = r.lock().unwrap();
        match data.deref().get(&un_locode) {
            Some(location) => Ok(location.clone()),
            None => Err(Box::new(ErrLocationNotFound)),
        }
    }

    fn find_all(&self) -> Result<Vec<Location>> {
        let r = self.0.clone();
        let data = r.lock().unwrap();
        let res = data.deref().values().map(|v| v.clone()).collect();
        Ok(res)
    }
}

#[derive(Debug)]
pub struct ErrLocationNotFound;

impl Error for ErrLocationNotFound {}

impl fmt::Display for ErrLocationNotFound {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Location is not found")
    }
}
