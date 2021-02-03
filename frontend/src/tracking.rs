use seed::{prelude::*, *};

struct_urls!();

pub struct Model;

pub fn view<Msg>(_model: &Model) -> Node<Msg> {
    div![h1!["Tracking ID"]]
}
