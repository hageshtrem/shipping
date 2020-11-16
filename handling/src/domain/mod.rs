pub mod handling;
pub mod location;

use std::hash::Hash;
use std::result::Result as StdResult;

pub type Result<T> = StdResult<T, Box<dyn std::error::Error>>;

pub trait Repository<K, V>
where
    K: Eq + Hash,
    V: Clone,
{
    fn store(&self, id: K, v: &V) -> Result<()>;
    fn find(&self, id: K) -> Result<V>;
    fn find_all(&self) -> Result<Vec<V>>;
}