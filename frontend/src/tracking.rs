use log::{error, info};
use seed::{prelude::*, *};
use serde::Deserialize;

pub fn init(_: Url, _: &mut impl Orders<Msg>) -> Model {
    Model {
        tracking_id: "".to_string(),
        cargo: None,
    }
}

pub struct Model {
    tracking_id: TrackingID,
    cargo: Option<Cargo>,
}

type TrackingID = String;

#[allow(non_snake_case)]
#[derive(Debug, Deserialize)]
pub struct Cargo {
    trackingId: String,
    statusText: String,
    origin: String,
    destination: String,
    eta: String,
    nextExpectedActivity: String,
    arrivalDeadline: String,
    isMisdirected: bool,
    events: Vec<Event>,
}

#[derive(Debug, Deserialize)]
pub struct Event {
    description: String,
    expected: bool,
}

pub enum Msg {
    Track,
    TrackingIDChanged(TrackingID),
    Fetched(fetch::Result<Cargo>),
}

pub fn update(msg: Msg, model: &mut Model, orders: &mut impl Orders<Msg>) {
    match msg {
        Msg::Track => {
            async fn send_message(id: TrackingID) -> fetch::Result<Cargo> {
                Request::new(format!("{}{}", crate::TRACKING_API_URL, id))
                    .method(Method::Get)
                    .fetch()
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            let id = model.tracking_id.clone();
            orders
                .skip()
                .perform_cmd(async { Msg::Fetched(send_message(id).await) });
        }
        Msg::TrackingIDChanged(id) => model.tracking_id = id,
        Msg::Fetched(Ok(cargo)) => {
            info!("{:?}", cargo);
            model.cargo = Some(cargo);
        }
        Msg::Fetched(Err(err)) => error!("{:?}", err),
    }
}

pub fn view(model: &Model) -> Node<Msg> {
    div![
        section![
            C!["section"],
            h1![C!["title"], "Tracking"],
            form![
                C!["form"],
                div![
                    C!["field", "has-addons"],
                    div![
                        C!["control"],
                        input![
                            C!["input"],
                            attrs! {
                                At::Value => model.tracking_id
                            },
                            input_ev(Ev::Input, Msg::TrackingIDChanged),
                        ]
                    ],
                    div![
                        C!["control"],
                        button![
                            C!["button"],
                            "Track!",
                            ev(Ev::Click, |event| {
                                event.prevent_default();
                                Msg::Track
                            })
                        ]
                    ]
                ]
            ]
        ],
        section![
            C!["section"],
            IF!(model.cargo.is_some() => view_cargo(model.cargo.as_ref().unwrap()))
        ]
    ]
}

fn view_cargo(cargo: &Cargo) -> Node<Msg> {
    div![
        C!["content"],
        p![format!(
            "Cargo {} is now: {}",
            cargo.trackingId, cargo.statusText
        )],
        if cargo.isMisdirected {
            vec![
                p![format!(
                    "Estimated time to arrival in {}: ?",
                    cargo.destination
                )],
                p!["Cargo is misdirected"],
            ]
        } else {
            vec![
                p![format!(
                    "Estimated time to arrival in {}: {}",
                    cargo.destination, cargo.eta
                )],
                p![format!("{}", cargo.nextExpectedActivity)],
            ]
        },
        p!["Delivery History:"],
        ul![cargo
            .events
            .iter()
            .map(|event| {
                li![
                    span![
                        C!["icon"],
                        if event.expected {
                            i![C!["fas", "fa-check-circle"]]
                        } else {
                            i![C!["fas", "fa-times-circle"]]
                        }
                    ],
                    *event.description
                ]
            })
            .collect::<Vec<Node<Msg>>>()]
    ]
}
