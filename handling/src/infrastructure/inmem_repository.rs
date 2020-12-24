use crate::domain::Repository;
use crate::Error;
use std::clone::Clone;
use std::collections::HashMap;
use std::hash::Hash;
use std::ops::{Deref, DerefMut};
use std::sync::{Arc, Mutex};

pub struct InmemRepository<K, V>(Arc<Mutex<HashMap<K, V>>>);

impl<K, V> InmemRepository<K, V> {
    pub fn new() -> Self {
        InmemRepository(Arc::new(Mutex::new(HashMap::new())))
    }
}

impl<K, V> Repository<K, V> for InmemRepository<K, V>
where
    K: Eq + Hash + std::fmt::Display + Send,
    V: Clone + Send,
{
    fn store(&self, key: K, value: &V) -> Result<(), Error> {
        let r = self.0.clone();
        let mut data = r.lock().unwrap();
        data.deref_mut().insert(key, value.clone());
        Ok(())
    }

    fn find(&self, key: K) -> Result<V, Error> {
        let r = self.0.clone();
        let data = r.lock().unwrap();
        match data.deref().get(&key) {
            Some(value) => Ok(value.clone()),
            None => Err(Error::RepositoryError(key.to_string())),
        }
    }

    fn find_all(&self) -> Result<Vec<V>, Error> {
        let r = self.0.clone();
        let data = r.lock().unwrap();
        let res = data.deref().values().map(|v| v.clone()).collect();
        Ok(res)
    }
}

impl<K, V> Clone for InmemRepository<K, V>
where
    K: Eq + Hash + std::fmt::Display,
    V: Clone,
{
    fn clone(&self) -> Self {
        InmemRepository(self.0.clone())
    }
}
