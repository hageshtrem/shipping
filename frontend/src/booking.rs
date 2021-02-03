use seed::{prelude::*, *};

pub struct Model;

struct_urls!();

pub fn view<Ms>(_model: &Model) -> Node<Ms> {
    div![h1!["Cargo Booking and Routing"],]
}
