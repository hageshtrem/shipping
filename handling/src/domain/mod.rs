pub mod handling;
pub mod location;

use std::hash::Hash;
use std::result::Result as StdResult;

pub type Result<T> = StdResult<T, Box<dyn std::error::Error>>;

pub trait Repository<K, V>: Clone + Send
where
    K: Eq + Hash + std::fmt::Display + Send,
    V: Clone + Send,
{
    fn store(&self, id: K, v: &V) -> Result<()>;
    fn find(&self, id: K) -> Result<V>;
    fn find_all(&self) -> Result<Vec<V>>;
}
