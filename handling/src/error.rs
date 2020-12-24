use std::{error, fmt};

#[derive(Debug, Clone)]
pub enum Error {
    InvalidArgument,
    HandlingError,
    RepositoryError(String),
    ParsingError,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::InvalidArgument => write!(f, "Provided argument is invalid"),
            Error::HandlingError => write!(f, "Event processing error"),
            Error::RepositoryError(msg) => write!(f, "Repository error: {}", msg),
            Error::ParsingError => write!(f, "Parsing error"),
        }
    }
}

impl error::Error for Error {}
