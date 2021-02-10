use log::{error, info};
use seed::{prelude::*, *};
use serde::{Deserialize, Serialize};
use serde_json::Value;

const ALL_CARGOS: &str = "cargo";
const NEW_CARGO: &str = "new";

pub fn init(url: Url, orders: &mut impl Orders<Msg>) -> Model {
    orders.subscribe(Msg::UrlChanged);
    orders.perform_cmd(async { Msg::CargosFetched(fetch_cargos().await) });
    Model {
        base_url: url.to_base_url(),
        page: Page::ListAllCargos,
        cargos: None,
        new_cargo: None,
    }
}

async fn fetch_cargos() -> fetch::Result<Data> {
    fetch(crate::BOOKING_API_URL)
        .await?
        .check_status()?
        .json()
        .await
}

pub struct Model {
    base_url: Url,
    page: Page,
    cargos: Option<Data>,
    new_cargo: Option<NewCargo>,
}

#[derive(Debug)]
enum Page {
    ListAllCargos,
    NewCargo,
}

impl Page {
    fn init(mut url: Url) -> Self {
        match url.remaining_path_parts().last() {
            Some(&NEW_CARGO) => Self::NewCargo,
            Some(&ALL_CARGOS) | None => Self::ListAllCargos,
            Some(_) => Self::ListAllCargos,
        }
    }
}

#[derive(Default, Deserialize)]
pub struct Data {
    cargos: Vec<Cargo>,
}

#[allow(non_snake_case)]
#[derive(Deserialize)]
pub struct Cargo {
    trackingId: String,
    origin: String,
    destination: String,
    routed: bool,
    misrouted: bool,
    arrivalDeadline: String,
    legs: Vec<Leg>,
}

#[allow(non_snake_case)]
#[derive(Deserialize)]
pub struct Leg {
    voyageNumber: String,
    loadLocation: String,
    loadTime: String,
    unloadLocation: String,
    unloadTime: String,
}

#[allow(non_snake_case)]
#[derive(Clone, Default, Serialize)]
pub struct NewCargo {
    origin: String,
    destination: String,
    deadline: String,
}

struct_urls!();
impl<'a> Urls<'a> {
    pub fn default(self) -> Url {
        self.all_cargos()
    }
    pub fn all_cargos(self) -> Url {
        self.base_url().add_path_part(ALL_CARGOS)
    }
    pub fn new_cargo(self) -> Url {
        self.base_url().add_path_part(NEW_CARGO)
    }
}

pub enum Msg {
    Book,
    NewCargoBooked(fetch::Result<Value>),
    UrlChanged(subs::UrlChanged),
    CargosFetched(fetch::Result<Data>),
    NewCargoOriginChanged(String),
    NewCargoDestinationChanged(String),
    NewCargoArrivalDeadlineChanged(String),
}

pub fn update(msg: Msg, model: &mut Model, orders: &mut impl Orders<Msg>) {
    match msg {
        Msg::Book => {
            async fn send_message(cargo: NewCargo) -> fetch::Result<Value> {
                Request::new(crate::BOOKING_API_URL)
                    .method(Method::Post)
                    .json(&cargo)?
                    .fetch()
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            if let Some(c) = &model.new_cargo {
                let mut c = c.clone();
                c.deadline = format!("{}:00Z", c.deadline);
                orders
                    .skip()
                    .perform_cmd(async { Msg::NewCargoBooked(send_message(c).await) });
            }
        }
        Msg::NewCargoBooked(res) => {
            info!("{:?}", res);
            orders.request_url(Urls::new(model.base_url.clone()).all_cargos());
        }
        Msg::UrlChanged(subs::UrlChanged(url)) => model.page = Page::init(url),
        Msg::CargosFetched(Ok(data)) => model.cargos = Some(data),
        Msg::CargosFetched(Err(err)) => error!("{:?}", err),
        Msg::NewCargoOriginChanged(origin) => match model.new_cargo.as_mut() {
            None => {
                model.new_cargo = Some(NewCargo {
                    origin: origin,
                    ..Default::default()
                });
            }
            Some(new_cargo) => (*new_cargo).origin = origin,
        },
        Msg::NewCargoDestinationChanged(destination) => match model.new_cargo.as_mut() {
            None => {
                model.new_cargo = Some(NewCargo {
                    destination: destination,
                    ..Default::default()
                });
            }
            Some(new_cargo) => (*new_cargo).destination = destination,
        },
        Msg::NewCargoArrivalDeadlineChanged(deadline) => match model.new_cargo.as_mut() {
            None => {
                model.new_cargo = Some(NewCargo {
                    deadline: deadline,
                    ..Default::default()
                });
            }
            Some(new_cargo) => (*new_cargo).deadline = deadline,
        },
    }
}

pub fn view(model: &Model, context: &crate::Context) -> Node<Msg> {
    let default = &Data { cargos: vec![] };
    let cargos = &model.cargos.as_ref().unwrap_or(default).cargos;
    div![
        h1!["Cargo Booking and Routing"],
        div![
            style! {
                St::Display => "inline",
            },
            a![
                attrs! { At::Href => Urls::new(model.base_url.clone()).all_cargos() },
                "List all cargos",
            ],
            "    ",
            a![
                attrs! { At::Href => Urls::new(model.base_url.clone()).new_cargo() },
                "Book new cargo",
            ],
        ],
        match model.page {
            Page::ListAllCargos => list_all_cargos_view(cargos),
            Page::NewCargo => new_cargo_view(model.new_cargo.as_ref(), context),
        },
    ]
}

fn list_all_cargos_view(cargos: &Vec<Cargo>) -> Node<Msg> {
    div![table![
        tr![
            th!["Tracking ID"],
            th!["Origin"],
            th!["Destination"],
            th!["Routed"]
        ],
        cargos.iter().map(|elem| {
            tr![
                td![elem.trackingId.clone()],
                td![elem.origin.clone()],
                td![elem.destination.clone()],
                td![elem.routed.to_string()],
            ]
        }),
    ]]
}

fn new_cargo_view(new_cargo_model: Option<&NewCargo>, context: &crate::Context) -> Node<Msg> {
    let default = &NewCargo::default();
    let new_cargo = new_cargo_model.unwrap_or(default);
    div![
        h2!["Book new cargo"],
        form![
            div![
                label!["Origin"],
                select![
                    option![attrs! {At::Value => AtValue::None}, "----"],
                    context
                        .locations
                        .iter()
                        .map(|loc| option![attrs! {At::Value => loc}, loc]),
                    input_ev(Ev::Input, Msg::NewCargoOriginChanged),
                ],
            ],
            div![
                label!["Destination"],
                select![
                    option![attrs! {At::Value => AtValue::None}, "----"],
                    context
                        .locations
                        .iter()
                        .map(|loc| option![attrs! {At::Value => loc}, loc]),
                    input_ev(Ev::Input, Msg::NewCargoDestinationChanged)
                ],
            ],
            div![
                label!["Arrival deadline"],
                input![
                    attrs! {
                        At::Type => "datetime-local",
                        At::Value => new_cargo.deadline,
                    },
                    input_ev(Ev::Input, Msg::NewCargoArrivalDeadlineChanged)
                ],
            ],
            div![button![
                "Book",
                ev(Ev::Click, |event| {
                    event.prevent_default();
                    Msg::Book
                })
            ],],
        ]
    ]
}
