fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(true)
        .build_client(true) //.out_dir("src/pb")
        .compile(
            &["../proto/handling.proto", "../proto/booking_events.proto"],
            &["../proto"],
        )?;
    Ok(())
}
