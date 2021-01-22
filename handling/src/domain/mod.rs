pub mod handling;
pub mod location;
pub mod voyage;

use crate::Error;
use std::hash::Hash;

pub trait Repository<K, V>: Clone + Send + Sync
where
    K: Eq + Hash + std::fmt::Display + Clone + Send,
    V: Clone + Send,
{
    fn store(&self, id: K, v: &V) -> Result<(), Error>;
    fn find(&self, id: K) -> Result<V, Error>;
    fn find_all(&self) -> Result<Vec<V>, Error>;
}
