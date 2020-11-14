pub mod handling;
pub mod location;

use std::result::Result as StdResult;

pub type Result<T> = StdResult<T, Box<dyn std::error::Error>>;

pub trait Repository<K, V> {
    fn store(&self, id: K, v: &V) -> Result<()>;
    fn find(&self, id: K) -> Result<V>;
    fn find_all(&self) -> Result<Vec<V>>;
}
