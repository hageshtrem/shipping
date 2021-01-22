use super::Repository;
use crate::Error;

pub type VoyageNumber = String;

#[derive(Clone)]
pub struct Voyage {}

// These voyages are hard-coded into the current pathfinder. Make sure
// they exist.
pub fn populate_repository<R: Repository<VoyageNumber, Voyage>>(
    repository: &R,
) -> Result<(), Error> {
    repository.store("0100S".to_string(), &Voyage {})?;
    repository.store("0200T".to_string(), &Voyage {})?;
    repository.store("0300A".to_string(), &Voyage {})?;
    repository.store("0301S".to_string(), &Voyage {})?;
    repository.store("0400S".to_string(), &Voyage {})?;
    Ok(())
}
