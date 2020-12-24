use handling::application::grpc_server::HandlingServiceImpl;
use handling::application::integration_events::event_handlers::{
    CargoDestinationChangedEventHandler, NewCargoBookedEventHandler,
};
use handling::application::pb::{CargoDestinationChanged, HandlingServiceServer, NewCargoBooked};
use handling::application::service::{EventHandlerImpl, ServiceImpl};
use handling::domain::handling::{Cargo, HandlingEventFactoryImpl, TrackingID, Voyage};
use handling::domain::{location, Repository};
use handling::infrastructure::inmem_repository::InmemRepository;
use handling::infrastructure::rabbitmq_eventbus::{EventBus, SubscribeManager};

use log::info;
use structopt::StructOpt;
use tonic::transport::Server;

/// Handling service
#[derive(StructOpt, Debug)]
struct Opt {
    /// Binding address
    #[structopt(long, env = "ADDR", default_value = "[::1]:5053")]
    addr: String,
    /// RabbitMQ
    #[structopt(
        long,
        env = "RABBIT_URI",
        default_value = "amqp://guest:guest@localhost:5672/%2f"
    )]
    rabbit_uri: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let opt = Opt::from_args();
    if std::env::var("RUST_LOG").is_err() {
        std::env::set_var("RUST_LOG", "info");
    }
    env_logger::init();

    // Dependencies
    let cargos = InmemRepository::new();
    let voyages = InmemRepository::new();
    voyages.store("v001".to_string(), &Voyage {}).unwrap();
    let locations = InmemRepository::new();
    location::store_sample_locations(&locations)?;
    let handling_events = InmemRepository::new();
    let unimp_event_handler = EventHandlerImpl {};
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
    let srv = ServiceImpl::new_service(handling_events, event_factory, unimp_event_handler);
    let addr = opt.addr.parse()?;
    let gservice = HandlingServiceImpl::new(srv);

    info!("Server started at {}", opt.addr);
    Server::builder()
        .add_service(HandlingServiceServer::new(gservice))
        .serve(addr)
        .await?;

    Ok(())
}
