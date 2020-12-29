use std::{error, fmt};

#[derive(Debug, Clone)]
pub enum Error {
    InvalidArgument,
    HandlingError,
    RepositoryError(String),
    ParsingError,
    EncodeError(prost::EncodeError),
    LapinError(lapin::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::InvalidArgument => write!(f, "Provided argument is invalid"),
            Error::HandlingError => write!(f, "Event processing error"),
            Error::RepositoryError(msg) => write!(f, "Repository error: {}", msg),
            Error::ParsingError => write!(f, "Parsing error"),
            Error::EncodeError(err) => write!(f, "{}", err.to_string()),
            Error::LapinError(err) => write!(f, "{}", err.to_string()),
        }
    }
}

impl error::Error for Error {}

impl From<prost::EncodeError> for Error {
    fn from(value: prost::EncodeError) -> Self {
        Error::EncodeError(value)
    }
}

impl From<lapin::Error> for Error {
    fn from(value: lapin::Error) -> Self {
        Error::LapinError(value)
    }
}
