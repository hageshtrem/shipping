use chrono::prelude::*;
use handling::application::pb::{
    HandlingEventType, HandlingServiceClient, RegisterHandlingEventRequest,
};
use prost_types::Timestamp;
use std::time::SystemTime;
use structopt::StructOpt;

#[derive(StructOpt, Debug)]
/// Handling service client
struct Opt {
    /// Binding address
    #[structopt(long, env = "ADDR", default_value = "127.0.0.1:5053")]
    addr: String,

    #[structopt(long, short)]
    id: String,

    #[structopt(long, short)]
    /// date in format dd.mm.yyyy
    completed: Option<String>,

    #[structopt(long, short)]
    voyage_number: String,

    #[structopt(long, short)]
    location: String,

    #[structopt(long, short)]
    event_type: HandlingEventType,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let opt = Opt::from_args();
    let completed = match opt.completed {
        Some(date) => {
            let naive_date = NaiveDate::parse_from_str(&date, "%d.%m.%Y")?;
            let naive_datetime: NaiveDateTime = naive_date.and_hms(0, 0, 0);
            let datetime_utc = DateTime::<Utc>::from_utc(naive_datetime, Utc);
            let ts: SystemTime = datetime_utc.into();
            Some(Timestamp::from(ts))
        }
        None => None,
    };

    let req = RegisterHandlingEventRequest {
        completed: completed,
        id: opt.id,
        voyage_number: opt.voyage_number,
        un_locode: opt.location,
        event_type: opt.event_type as i32,
    };
    let mut client = HandlingServiceClient::connect(format!("http://{}", opt.addr)).await?;
    let resp = client.register_handling_event(req).await?;
    println!("{:#?}", resp);

    Ok(())
}
