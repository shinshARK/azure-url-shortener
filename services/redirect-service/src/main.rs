use axum::{routing::get, Router};
use std::env;
use std::net::SocketAddr;

// Register modules
mod db;
mod analytics;
mod models;
mod handlers;

#[tokio::main]
async fn main() {
    let port_key = "FUNCTIONS_CUSTOMHANDLER_PORT";
    let port: u16 = env::var(port_key)
        .unwrap_or_else(|_| "3000".to_string())
        .parse()
        .expect("Port must be a number");

    let app = Router::new().route("/redirect/{short_code}", get(handlers::handle_redirect));

    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    println!("Listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
