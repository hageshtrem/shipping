use log::{error, info};
use seed::{prelude::*, *};
use serde::Serialize;

pub fn init(_: Url, _: &mut impl Orders<Msg>) -> Model {
    Model {
        event: Event::default(),
    }
}

pub struct Model {
    event: Event,
}

#[derive(Default, Serialize, Clone)]
pub struct Event {
    completed: String,
    id: String,
    voyage_number: String,
    un_locode: String,
    event_type: String,
}

pub enum Msg {
    Register,
    Fetched(fetch::Result<()>),
    DateTimeChanged(String),
    TrackingIDChanged(String),
    VoyageCanged(String),
    LocationChanged(String),
    EventTypeChanged(String),
}

pub fn update(msg: Msg, model: &mut Model, orders: &mut impl Orders<Msg>) {
    match msg {
        Msg::Register => {
            async fn send_message(event: Event) -> fetch::Result<()> {
                Request::new(format!("{}{}", crate::HANDLING_API_URL, event.id))
                    .method(Method::Post)
                    .json(&event)?
                    .fetch()
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            let mut e = model.event.clone();
            e.completed = format!("{}:00Z", e.completed);
            orders
                .skip()
                .perform_cmd(async { Msg::Fetched(send_message(e).await) });
        }
        Msg::Fetched(Ok(_)) => info!("Registered!"),
        Msg::Fetched(Err(err)) => error!("{:?}", err),
        Msg::DateTimeChanged(datetime) => model.event.completed = datetime,
        Msg::TrackingIDChanged(id) => model.event.id = id,
        Msg::VoyageCanged(voyage_number) => model.event.voyage_number = voyage_number,
        Msg::LocationChanged(un_locode) => model.event.un_locode = un_locode,
        Msg::EventTypeChanged(event_type) => model.event.event_type = event_type,
    }
}

pub fn view(model: &Model) -> Node<Msg> {
    let locations = vec![
        "SESTO", "SEGOT", "AUMEL", "CNHKG", "CNSHA", "CNHGH", "USNYC", "USCHI", "USDAL", "JNTKO",
        "DEHAM", "NLRTM", "FIHEL",
    ];
    let voyages = vec!["0100S", "0200T", "0300A", "0301S", "0400S"];
    div![
        h1!["Incident Logging Application"],
        form![
            label!["Time: "],
            input![
                attrs! {
                    At::Type => "datetime-local",
                    At::Value => model.event.completed,
                },
                input_ev(Ev::Input, Msg::DateTimeChanged)
            ],
            label!["Tracking ID: "],
            input![
                attrs! {
                    At::Value => model.event.id,
                },
                input_ev(Ev::Input, Msg::TrackingIDChanged)
            ],
            label!["Voyage: "],
            select![
                option![attrs! {At::Value => AtValue::None}, "----"],
                voyages
                    .iter()
                    .map(|code| option![attrs! {At::Value => code}, code]),
                input_ev(Ev::Input, Msg::VoyageCanged)
            ],
            label!["Location: "],
            select![
                option![attrs! {At::Value => AtValue::None}, "----"],
                locations
                    .iter()
                    .map(|loc| option![attrs! {At::Value => loc}, loc]),
                input_ev(Ev::Input, Msg::LocationChanged)
            ],
            label!["Event Type: "],
            select![
                option![attrs! {At::Value => AtValue::None}, "----"],
                option![attrs! {At::Value => "NotHandled"}, "NotHandled"],
                option![attrs! {At::Value => "Receive"}, "Receive"],
                option![attrs! {At::Value => "Load"}, "Load"],
                option![attrs! {At::Value => "Unload"}, "Unload"],
                option![attrs! {At::Value => "Claim"}, "Claim"],
                option![attrs! {At::Value => "Customs"}, "Customs"],
                input_ev(Ev::Input, Msg::EventTypeChanged),
            ],
            button![
                "Register",
                ev(Ev::Click, |event| {
                    event.prevent_default();
                    Msg::Register
                })
            ],
        ]
    ]
}
