use std::{error::Error, fmt};

#[derive(Debug)]
pub struct ErrInvalidArgument;

impl Error for ErrInvalidArgument {}

impl fmt::Display for ErrInvalidArgument {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Provided argument is invalid")
    }
}
