#![allow(dead_code)]

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
