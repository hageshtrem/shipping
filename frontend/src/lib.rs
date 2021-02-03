#![allow(clippy::must_use_candidate)]

use seed::{prelude::*, *};

mod booking;
mod handling;
mod tracking;

const BOOKING: &str = "booking";
const TRACKING: &str = "tracking";
const HANDLING: &str = "handling";

// ------ ------
//     Init
// ------ ------

fn init(url: Url, orders: &mut impl Orders<Msg>) -> Model {
    orders.subscribe(Msg::UrlChanged);
    Model {
        ctx: Context {
            tracking_id: "".to_string(),
        },
        base_url: url.to_base_url(),
        page: Page::Booking,
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
    pub tracking_id: TrackingID,
}

type TrackingID = String;

// ------ Page ------

enum Page {
    Booking,
    Tracking,
    Handling,
}

impl Page {
    fn init(mut url: Url) -> Self {
        match url.next_path_part() {
            Some(TRACKING) => Self::Tracking,
            Some(HANDLING) => Self::Handling,
            Some(BOOKING) | _ => Self::Booking,
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
    pub fn tracking(self) -> tracking::Urls<'a> {
        tracking::Urls::new(self.base_url().add_path_part(TRACKING))
    }
    pub fn handling(self) -> handling::Urls<'a> {
        handling::Urls::new(self.base_url().add_path_part(HANDLING))
    }
}

// ------ ------
//    Update
// ------ ------

enum Msg {
    UrlChanged(subs::UrlChanged),
}

fn update(msg: Msg, model: &mut Model, _: &mut impl Orders<Msg>) {
    match msg {
        Msg::UrlChanged(subs::UrlChanged(mut url)) => {
            model.page = match url.next_path_part() {
                None | Some(BOOKING) => Page::Booking,
                Some(TRACKING) => Page::Tracking,
                Some(HANDLING) => Page::Handling,
                Some(_) => Page::Booking,
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
        match &model.page {
            Page::Booking => booking::view(&booking::Model {}),
            Page::Tracking => tracking::view(&tracking::Model {}),
            Page::Handling => handling::view(&handling::Model {}),
        },
    ]
}

fn header(base_url: &Url) -> Node<Msg> {
    nav![
        a![
            attrs! { At::Href => Urls::new(base_url).booking().base_url() },
            "Booking",
        ],
        " | ",
        a![
            attrs! { At::Href => Urls::new(base_url).tracking().base_url() },
            "Tracking",
        ],
        " | ",
        a![
            attrs! { At::Href => Urls::new(base_url).handling().base_url() },
            "Handling",
        ],
    ]
}

// ------ ------
//     Start
// ------ ------

#[wasm_bindgen(start)]
pub fn start() {
    App::start("app", init, update, view);
}
