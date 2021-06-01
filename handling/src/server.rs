use handling::application::grpc_server::{HandlingServiceImpl, NamedHandlingServiceImpl};
use handling::application::integration_events::{
    CargoDestinationChangedEventHandler, NewCargoBookedEventHandler,
};
use handling::application::logging_service::LoggingService;
use handling::application::pb::{CargoDestinationChanged, HandlingServiceServer, NewCargoBooked};
use handling::application::service::ServiceImpl;
use handling::domain::handling::{Cargo, HandlingEventFactoryImpl, TrackingID};
use handling::domain::{location, voyage};
use handling::infrastructure::inmem_repository::InmemRepository;
use handling::infrastructure::rabbitmq_eventbus::{EventBus, SubscribeManager};

use chrono::prelude::*;
use log::{info, LevelFilter};
use log4rs::{
    append::{console::ConsoleAppender, file::FileAppender},
    config::{Appender, Config, Root},
    encode::json::JsonEncoder,
};
use std::fs;
use std::path::Path;
use structopt::StructOpt;
use tonic::transport::Server;

/// Handling service
#[derive(StructOpt, Debug)]
struct Opt {
    /// Binding address
    #[structopt(long, env = "ADDR", default_value = "127.0.0.1:5053")]
    addr: String,
    /// RabbitMQ
    #[structopt(
        long,
        env = "RABBIT_URI",
        default_value = "amqp://guest:guest@127.0.0.1:5672/%2f"
    )]
    rabbit_uri: String,
    /// Directory for logs
    #[structopt(long, env = "LOG_DIR", default_value = "/var/log/handling")]
    log_dir: String,
}

fn init_logger(dir: &str) -> Result<(), Box<dyn std::error::Error>> {
    fs::create_dir_all(dir)?;
    let file_path =
        Path::new(dir).join(format!("handling-{}.log", Utc::today().format("%m.%d.%Y")));

    let stdout = ConsoleAppender::builder()
        .encoder(Box::new(JsonEncoder::new()))
        .build();

    let logfile = FileAppender::builder()
        .encoder(Box::new(JsonEncoder::new()))
        .build(file_path)?;

    let config = Config::builder()
        .appender(Appender::builder().build("logfile", Box::new(logfile)))
        .appender(Appender::builder().build("stdout", Box::new(stdout)))
        .build(
            Root::builder()
                .appender("logfile")
                .appender("stdout")
                .build(LevelFilter::Info),
        )?;

    log4rs::init_config(config)?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let opt = Opt::from_args();
    init_logger(&opt.log_dir)?;

    // Dependencies
    let cargos = InmemRepository::new();
    let voyages = InmemRepository::new();
    voyage::populate_repository(&voyages)?;
    let locations = InmemRepository::new();
    location::store_sample_locations(&locations)?;
    let handling_events = InmemRepository::new();
    let event_factory = HandlingEventFactoryImpl::new(cargos.clone(), voyages, locations);

    // IntegrationEventBus
    let new_cargo_eh = NewCargoBookedEventHandler::new(cargos.clone());
    let cargo_dest_changed_eh = CargoDestinationChangedEventHandler::new(cargos);
    let mut event_bus = EventBus::new(&opt.rabbit_uri).await?;
    event_bus.subscribe::<NewCargoBooked, NewCargoBookedEventHandler<InmemRepository<TrackingID, Cargo>>>(
            new_cargo_eh,
        ).await?;
    event_bus.subscribe::<CargoDestinationChanged, CargoDestinationChangedEventHandler<InmemRepository<TrackingID, Cargo>>>(
            cargo_dest_changed_eh
        ).await?;

    // Service
    let srv = ServiceImpl::new_service(handling_events, event_factory, event_bus);
    let srv = LoggingService::new(Box::new(srv));
    let addr = opt.addr.parse()?;
    let gservice = HandlingServiceImpl::new(srv);

    // Health
    let (mut health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter
        .set_serving::<NamedHandlingServiceImpl>()
        .await;

    info!("Server started at {}", opt.addr);
    Server::builder()
        .add_service(health_service)
        .add_service(HandlingServiceServer::new(gservice))
        .serve(addr)
        .await?;

    Ok(())
}
