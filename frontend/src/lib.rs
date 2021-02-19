#![allow(clippy::wildcard_imports)]

use seed::{prelude::*, *};

mod booking;
mod handling;
mod tracking;

const BOOKING: &str = "booking";
const TRACKING: &str = "tracking";
const HANDLING: &str = "handling";

pub(crate) const BOOKING_API_URL: &str = "http://localhost:8080/booking/v1/cargos/";
pub(crate) const TRACKING_API_URL: &str = "http://localhost:8080/tracking/v1/cargos/";
pub(crate) const HANDLING_API_URL: &str = "http://localhost:8080/handling/v1/cargos/";

// ------ ------
//     Init
// ------ ------

fn init(url: Url, orders: &mut impl Orders<Msg>) -> Model {
    orders.subscribe(Msg::UrlChanged);
    Model {
        ctx: Context {
            tracking_id: None,
            locations: vec![
                "SESTO", "SEGOT", "AUMEL", "CNHKG", "CNSHA", "CNHGH", "USNYC", "USCHI", "USDAL",
                "JNTKO", "DEHAM", "NLRTM", "FIHEL",
            ],
            voyages: vec!["0100S", "0200T", "0300A", "0301S", "0400S"],
        },
        base_url: url.to_base_url(),
        page: Page::init(url, orders),
    }
}

// ------ ------
//     Model
// ------ ------

struct Model {
    ctx: Context,
    base_url: Url,
    page: Page,
}

// ------ Context ------

pub struct Context {
    pub tracking_id: Option<String>,
    pub locations: Vec<&'static str>,
    pub voyages: Vec<&'static str>,
}

// ------ Page ------

enum Page {
    Booking(booking::Model),
    Tracking(tracking::Model),
    Handling(handling::Model),
}

impl Page {
    fn init(mut url: Url, orders: &mut impl Orders<Msg>) -> Self {
        match url.next_path_part() {
            Some(TRACKING) => {
                Self::Tracking(tracking::init(url, &mut orders.proxy(Msg::TrackingMsg)))
            }
            Some(HANDLING) => {
                Self::Handling(handling::init(url, &mut orders.proxy(Msg::HandlingMsg)))
            }
            Some(BOOKING) | _ => {
                Self::Booking(booking::init(url, &mut orders.proxy(Msg::BookingMsg)))
            }
        }
    }
}

// ------ ------
//     Urls
// ------ ------

struct_urls!();
impl<'a> Urls<'a> {
    pub fn home(self) -> Url {
        self.base_url()
    }
    pub fn booking(self) -> booking::Urls<'a> {
        booking::Urls::new(self.base_url().add_path_part(BOOKING))
    }
    pub fn tracking(self) -> Url {
        self.base_url().add_path_part(TRACKING)
    }
    pub fn handling(self) -> Url {
        self.base_url().add_path_part(HANDLING)
    }
}

// ------ ------
//    Update
// ------ ------

enum Msg {
    UrlChanged(subs::UrlChanged),
    BookingMsg(booking::Msg),
    TrackingMsg(tracking::Msg),
    HandlingMsg(handling::Msg),
}

fn update(msg: Msg, model: &mut Model, orders: &mut impl Orders<Msg>) {
    match msg {
        Msg::UrlChanged(subs::UrlChanged(url)) => model.page = Page::init(url, orders),
        Msg::BookingMsg(msg) => {
            if let Page::Booking(model) = &mut model.page {
                booking::update(msg, model, &mut orders.proxy(Msg::BookingMsg))
            }
        }
        Msg::TrackingMsg(msg) => {
            if let Page::Tracking(model) = &mut model.page {
                tracking::update(msg, model, &mut orders.proxy(Msg::TrackingMsg))
            }
        }
        Msg::HandlingMsg(msg) => {
            if let Page::Handling(model) = &mut model.page {
                handling::update(msg, model, &mut orders.proxy(Msg::HandlingMsg))
            }
        }
    }
}

// ------ ------
//     View
// ------ ------

fn view(model: &Model) -> impl IntoNodes<Msg> {
    vec![
        header(&model.base_url),
        div![
            C!["container"],
            match &model.page {
                Page::Booking(booking_model) => {
                    booking::view(booking_model, &model.ctx).map_msg(Msg::BookingMsg)
                }
                Page::Tracking(tracking_model) => {
                    tracking::view(tracking_model).map_msg(Msg::TrackingMsg)
                }
                Page::Handling(handling_model) => {
                    handling::view(handling_model, &model.ctx).map_msg(Msg::HandlingMsg)
                }
            },
        ],
    ]
}

fn header(base_url: &Url) -> Node<Msg> {
    nav![div![
        C!["navbar-menu"],
        div![
            C!["navbar-start"],
            a![
                C!["navbar-item"],
                attrs! { At::Href => Urls::new(base_url).booking().default() },
                "Booking",
            ],
            a![
                C!["navbar-item"],
                attrs! { At::Href => Urls::new(base_url).tracking() },
                "Tracking",
            ],
            a![
                C!["navbar-item"],
                attrs! { At::Href => Urls::new(base_url).handling() },
                "Handling",
            ]
        ]
    ]]
}

// ------ ------
//     Start
// ------ ------

#[wasm_bindgen(start)]
pub fn start() {
    console_log::init().expect("error initializing logger");
    App::start("app", init, update, view);
}
