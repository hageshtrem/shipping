use super::Repository;
use crate::Error;

// UNLocode is the United Nations location code that uniquely identifies a
// particular location.
//
// http://www.unece.org/cefact/locode/
// http://www.unece.org/cefact/locode/DocColumnDescription.htm#LOCODE
pub type UNLocode = String;

// Location is a location is our model is stops on a journey, such as cargo
// origin or destination, or carrier movement endpoints.
#[derive(Clone)]
pub struct Location {
    un_locode: UNLocode,
    name: String,
}

#[allow(non_snake_case)]
pub fn store_sample_locations<R: Repository<UNLocode, Location>>(
    repository: &R,
) -> Result<(), Error> {
    let SESTO = &Location {
        un_locode: "SESTO".to_string(),
        name: "Stockholm".to_string(),
    };
    let AUMEL = &Location {
        un_locode: "AUMEL".to_string(),
        name: "Melbourne".to_string(),
    };
    let CNHKG = &Location {
        un_locode: "CNHKG".to_string(),
        name: "Hongkong".to_string(),
    };
    let USNYC = &Location {
        un_locode: "USNYC".to_string(),
        name: "New York".to_string(),
    };
    let USCHI = &Location {
        un_locode: "USCHI".to_string(),
        name: "Chicago".to_string(),
    };
    let JNTKO = &Location {
        un_locode: "JNTKO".to_string(),
        name: "Tokyo".to_string(),
    };
    let DEHAM = &Location {
        un_locode: "DEHAM".to_string(),
        name: "Hamburg".to_string(),
    };
    let NLRTM = &Location {
        un_locode: "NLRTM".to_string(),
        name: "Rotterdam".to_string(),
    };
    let FIHEL = &Location {
        un_locode: "FIHEL".to_string(),
        name: "Helsinki".to_string(),
    };
    repository.store(SESTO.un_locode.clone(), SESTO)?;
    repository.store(AUMEL.un_locode.clone(), AUMEL)?;
    repository.store(CNHKG.un_locode.clone(), CNHKG)?;
    repository.store(USNYC.un_locode.clone(), USNYC)?;
    repository.store(USCHI.un_locode.clone(), USCHI)?;
    repository.store(JNTKO.un_locode.clone(), JNTKO)?;
    repository.store(DEHAM.un_locode.clone(), DEHAM)?;
    repository.store(NLRTM.un_locode.clone(), NLRTM)?;
    repository.store(FIHEL.un_locode.clone(), FIHEL)?;

    Ok(())
}
