use log::{error, info};
use seed::{prelude::*, *};
use serde::{Deserialize, Serialize};
use serde_json::Value;

const ALL_CARGOS: &str = "cargo";
const CARGO_DETAILS: &str = "details";
const ROUTE_CARGO: &str = "route";
const CHANGE_DESTINATION: &str = "destination";
const NEW_CARGO: &str = "new";

pub fn init(url: Url, orders: &mut impl Orders<Msg>) -> Model {
    orders.subscribe(Msg::UrlChanged);
    orders.perform_cmd(async { Msg::CargosFetched(fetch_cargos().await) });
    Model {
        base_url: url.to_base_url(),
        page: Page::ListAllCargos,
        cargos: None,
        new_cargo: None,
        itineraries: None,
        change_destination_request: None,
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
    itineraries: Option<Itineraries>,
    change_destination_request: Option<ChangeDestinationRequest>,
}

#[derive(Debug)]
enum Page {
    ListAllCargos,
    Details(String),
    Route(String),
    ChangeDestination(String),
    NewCargo,
}

impl Page {
    fn init(mut url: Url, orders: &mut impl Orders<Msg>) -> Self {
        match url.remaining_path_parts().as_slice() {
            [.., NEW_CARGO] => Self::NewCargo,
            [.., CARGO_DETAILS, id @ _] => Self::Details(id.to_string()),
            [.., ROUTE_CARGO, id @ _] => {
                orders.send_msg(Msg::GetRoutes(id.to_string()));
                Self::Route(id.to_string())
            }
            [.., CHANGE_DESTINATION, id @ _] => Self::ChangeDestination(id.to_string()),
            [.., ALL_CARGOS] => Self::ListAllCargos,
            [] | [..] => Self::ListAllCargos,
        }
    }
}

#[derive(Default, Deserialize)]
pub struct Data {
    cargos: Vec<Cargo>,
}

impl Data {
    fn get_cargo_by_id(&self, id: &str) -> Option<&Cargo> {
        self.cargos.iter().by_ref().find(|&c| c.trackingId == *id)
    }
}

#[allow(non_snake_case)]
#[derive(Debug, Deserialize)]
pub struct Cargo {
    trackingId: TrackingID,
    origin: String,
    destination: String,
    routed: bool,
    misrouted: bool,
    arrivalDeadline: String,
    legs: Vec<Leg>,
}

type TrackingID = String;

#[derive(Clone, Debug, Deserialize)]
pub struct Itineraries(Vec<Itinerary>);

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Itinerary {
    legs: Vec<Leg>,
}

#[allow(non_snake_case)]
#[derive(Clone, Debug, Deserialize, Serialize)]
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

#[allow(non_snake_case)]
#[derive(Clone, Default, Serialize)]
pub struct ChangeDestinationRequest {
    trackingId: TrackingID,
    destination: String,
}

struct_urls!();
impl<'a> Urls<'a> {
    pub fn default(self) -> Url {
        self.all_cargos()
    }
    pub fn all_cargos(self) -> Url {
        self.base_url().add_path_part(ALL_CARGOS)
    }
    pub fn details(self, id: &str) -> Url {
        self.base_url()
            .add_path_part(CARGO_DETAILS)
            .add_path_part(id)
    }
    pub fn route(self, id: &str) -> Url {
        self.base_url().add_path_part(ROUTE_CARGO).add_path_part(id)
    }
    pub fn destination(self, id: &str) -> Url {
        self.base_url()
            .add_path_part(CHANGE_DESTINATION)
            .add_path_part(id)
    }
    pub fn new_cargo(self) -> Url {
        self.base_url().add_path_part(NEW_CARGO)
    }
}

pub enum Msg {
    UrlChanged(subs::UrlChanged),
    CargosFetched(fetch::Result<Data>),
    Book,
    NewCargoBooked(fetch::Result<Value>),
    GetRoutes(TrackingID),
    RoutesFetched(fetch::Result<Itineraries>),
    AssignCargoToRoute(TrackingID, Itinerary),
    CargoToRouteAssigned(fetch::Result<Value>),
    ChangeDestination(TrackingID),
    DestinationChanged(fetch::Result<Value>),
    ChangeDestinationRequestDestinationChanged(String),
    NewCargoOriginChanged(String),
    NewCargoDestinationChanged(String),
    NewCargoArrivalDeadlineChanged(String),
}

pub fn update(msg: Msg, model: &mut Model, orders: &mut impl Orders<Msg>) {
    match msg {
        Msg::UrlChanged(subs::UrlChanged(url)) => model.page = Page::init(url, orders),
        Msg::CargosFetched(Ok(data)) => model.cargos = Some(data),
        Msg::CargosFetched(Err(err)) => error!("{:?}", err),
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
        Msg::GetRoutes(id) => {
            async fn send_message(id: String) -> fetch::Result<Itineraries> {
                fetch(format!("{}{}/routes", crate::BOOKING_API_URL, id))
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            orders
                .skip()
                .perform_cmd(async { Msg::RoutesFetched(send_message(id).await) });
        }
        Msg::RoutesFetched(Ok(itineraries)) => {
            info!("Routes fetched!"); // TODO: two times
            model.itineraries = Some(itineraries);
        }
        Msg::RoutesFetched(Err(err)) => error!("{:?}", err),
        Msg::AssignCargoToRoute(id, itinerary) => {
            info!("Assign! {} \n {:?}", id, itinerary);
            async fn send_message(id: String, itinerary: Itinerary) -> fetch::Result<Value> {
                Request::new(format!("{}{}/route", crate::BOOKING_API_URL, id))
                    .method(Method::Put)
                    .json(&itinerary)?
                    .fetch()
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            orders.skip().perform_cmd(async {
                Msg::CargoToRouteAssigned(send_message(id, itinerary).await)
            });
        }
        Msg::CargoToRouteAssigned(Ok(value)) => {
            info!("{}", value);
            orders.request_url(Urls::new(model.base_url.clone()).all_cargos());
        }
        Msg::CargoToRouteAssigned(Err(err)) => error!("{:?}", err),
        Msg::ChangeDestination(id) => {
            info!("Change destination: {}", id);
            async fn send_message(req: ChangeDestinationRequest) -> fetch::Result<Value> {
                Request::new(format!("{}{}", crate::BOOKING_API_URL, req.trackingId))
                    .method(Method::Put)
                    .json(&req)?
                    .fetch()
                    .await?
                    .check_status()?
                    .json()
                    .await
            }
            if let Some(req) = &model.change_destination_request {
                let mut req = req.clone();
                req.trackingId = id;
                orders
                    .skip()
                    .perform_cmd(async { Msg::DestinationChanged(send_message(req).await) });
            }
        }
        Msg::DestinationChanged(Ok(value)) => {
            info!("{}", value);
            orders.request_url(Urls::new(model.base_url.clone()).all_cargos());
        }
        Msg::DestinationChanged(Err(err)) => error!("{:?}", err),
        Msg::ChangeDestinationRequestDestinationChanged(dst) => {
            match model.change_destination_request.as_mut() {
                None => {
                    model.change_destination_request = Some(ChangeDestinationRequest {
                        destination: dst,
                        ..Default::default()
                    })
                }
                Some(req) => (*req).destination = dst,
            }
        }
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
        h1![C!["title"], "Cargo Booking and Routing"],
        div![
            C!["columns"],
            div![
                C!["column", "is-2"],
                a![
                    attrs! { At::Href => Urls::new(model.base_url.clone()).all_cargos() },
                    "List all cargos",
                ]
            ],
            div![
                C!["column", "is-2"],
                a![
                    attrs! { At::Href => Urls::new(model.base_url.clone()).new_cargo() },
                    "Book new cargo",
                ]
            ]
        ],
        match &model.page {
            Page::ListAllCargos => list_all_cargos_view(&model.base_url, cargos),
            Page::Details(id) => {
                if let Some(cargo) = model
                    .cargos
                    .as_ref()
                    .and_then(|data| data.get_cargo_by_id(id))
                {
                    details_of_cargo_view(&model.base_url, cargo)
                } else {
                    empty![]
                }
            }
            Page::Route(id) => route_cargo_view(model, id),
            Page::ChangeDestination(id) => change_destination_view(id, context),
            Page::NewCargo => new_cargo_view(model.new_cargo.as_ref(), context),
        },
    ]
}

fn list_all_cargos_view(base_url: &Url, cargos: &Vec<Cargo>) -> Node<Msg> {
    table![
        C!["table"],
        thead![tr![
            th!["Tracking ID"],
            th!["Origin"],
            th!["Destination"],
            th!["Routed"]
        ]],
        tbody![cargos.iter().map(|elem| {
            tr![
                td![a![
                    attrs! { At::Href => Urls::new(base_url.clone()).details(&elem.trackingId) },
                    elem.trackingId.clone(),
                ],],
                td![elem.origin.clone()],
                td![elem.destination.clone()],
                td![span![
                    C!["icon"],
                    if elem.routed {
                        i![C!["fas", "fa-check-circle"]]
                    } else {
                        i![C!["fas", "fa-circle"]]
                    }
                ]]
            ]
        })]
    ]
}

fn new_cargo_view(new_cargo_model: Option<&NewCargo>, context: &crate::Context) -> Node<Msg> {
    let default = &NewCargo::default();
    let new_cargo_model = new_cargo_model.unwrap_or(default);
    div![
        h2![C!["subtitle"], "Book new cargo"],
        div![
            C!["column", "is-4"],
            form![
                div![
                    C!["field"],
                    label![C!["label"], "Origin"],
                    div![
                        C!["control", "is-expanded"],
                        div![
                            C!["select", "is-fullwidth"],
                            select![
                                option![attrs! {At::Value => AtValue::None}, "----"],
                                context
                                    .locations
                                    .iter()
                                    .map(|loc| option![attrs! {At::Value => loc}, loc]),
                                input_ev(Ev::Input, Msg::NewCargoOriginChanged),
                            ]
                        ]
                    ]
                ],
                div![
                    C!["field"],
                    label![C!["label"], "Destination"],
                    div![
                        C!["control", "is-expanded"],
                        div![
                            C!["select", "is-fullwidth"],
                            select![
                                option![attrs! {At::Value => AtValue::None}, "----"],
                                context
                                    .locations
                                    .iter()
                                    .map(|loc| option![attrs! {At::Value => loc}, loc]),
                                input_ev(Ev::Input, Msg::NewCargoDestinationChanged)
                            ]
                        ]
                    ]
                ],
                div![
                    C!["field"],
                    label![C!["label"], "Arrival deadline"],
                    div![
                        C!["control", "is-expanded"],
                        input![
                            C!["input"],
                            attrs! {
                                At::Type => "datetime-local",
                                At::Value => new_cargo_model.deadline,
                            },
                            input_ev(Ev::Input, Msg::NewCargoArrivalDeadlineChanged)
                        ]
                    ]
                ],
                div![
                    C!["field"],
                    div![
                        C!["control"],
                        button![
                            C!["button"],
                            "Book",
                            ev(Ev::Click, |event| {
                                event.prevent_default();
                                Msg::Book
                            })
                        ]
                    ]
                ]
            ]
        ]
    ]
}

fn details_of_cargo_view(base_url: &Url, cargo: &Cargo) -> Node<Msg> {
    div![
        h2![
            C!["subtitle"],
            format!("Details for cargo {}", cargo.trackingId)
        ],
        table![
            C!["table"],
            tbody![
                tr![td!["Origin"], td![*cargo.origin]],
                tr![
                    td!["Destination"],
                    td![
                        div![*cargo.destination],
                        div![a![
                            attrs! { At::Href => Urls::new(base_url.clone()).destination(&cargo.trackingId) },
                            "Change destination"
                        ]]
                    ]
                ],
                tr![td!["Arrival deadline"], td![*cargo.arrivalDeadline]]
            ]
        ],
        div![if cargo.routed {
            itinerary_view(cargo.legs.as_ref())
        } else {
            div![
                "Not routed - ",
                a![
                    attrs! { At::Href => Urls::new(base_url.clone()).route(&cargo.trackingId) },
                    "Route this cargo"
                ]
            ]
        }]
    ]
}

fn route_cargo_view(model: &Model, id: &str) -> Node<Msg> {
    div![
        h2![C!["subtitle"], "Select route"],
        if let Some(cargo) = model
            .cargos
            .as_ref()
            .and_then(|data| data.get_cargo_by_id(id))
        {
            div![
                format!(
                    "Cargo {} is going from {} to {}",
                    cargo.trackingId, cargo.origin, cargo.destination
                ),
                if let Some(itineraries) = &model.itineraries {
                    div![itineraries.0.iter().map(|itinerary| itinerary.clone()).map(
                        move |itinerary| {
                            div![itinerary_view(itinerary.legs.as_ref()), {
                                let id = id.to_string();
                                button![
                                    C!["button"],
                                    "Assign cargo to this route",
                                    ev(Ev::Click, move |_| Msg::AssignCargoToRoute(id, itinerary))
                                ]
                            }]
                        }
                    )]
                } else {
                    empty![]
                }
            ]
        } else {
            empty![]
        }
    ]
}

fn change_destination_view(id: &str, context: &crate::Context) -> Node<Msg> {
    div![
        h2![
            C!["subtitle"],
            format!("Change destination of cargo {}", id)
        ],
        div![
            C!["field", "has-addons"],
            div![
                C!["control"],
                div![
                    C!["select"],
                    select![
                        option![attrs! {At::Value => AtValue::None}, "----"],
                        context
                            .locations
                            .iter()
                            .map(|loc| option![attrs! {At::Value => loc}, loc]),
                        input_ev(Ev::Input, Msg::ChangeDestinationRequestDestinationChanged)
                    ]
                ]
            ],
            div![
                C!["control"],
                button![C!["button"], "Change", {
                    let id = id.to_string();
                    ev(Ev::Click, move |event| {
                        event.prevent_default();
                        Msg::ChangeDestination(id)
                    })
                }]
            ]
        ]
    ]
}

fn itinerary_view(legs: &Vec<Leg>) -> Node<Msg> {
    table![
        C!["table"],
        thead![tr![
            th!["Voyage number"],
            th!["Load"],
            th!["Load Time"],
            th!["Unload"],
            th!["Unload Time"],
        ]],
        tbody![legs.iter().map(|leg| {
            tr![
                td![*leg.voyageNumber],
                td![*leg.loadLocation],
                td![*leg.loadTime],
                td![*leg.unloadLocation],
                td![*leg.unloadTime],
            ]
        })]
    ]
}
