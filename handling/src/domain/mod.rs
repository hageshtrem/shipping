pub mod handling;
pub mod location;

use crate::Error;
use std::hash::Hash;

pub trait Repository<K, V>: Clone + Send
where
    K: Eq + Hash + std::fmt::Display + Send,
    V: Clone + Send,
{
    fn store(&self, id: K, v: &V) -> Result<(), Error>;
    fn find(&self, id: K) -> Result<V, Error>;
    fn find_all(&self) -> Result<Vec<V>, Error>;
}
