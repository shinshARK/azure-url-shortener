use axum::{
    extract::Request,
    routing::get,
    Router,
    middleware::{self, Next},
    response::Response,
};
use std::env;
use std::net::SocketAddr;

mod db;
mod analytics;
mod models;
mod handlers;

// New Debug Logger
async fn log_request(req: Request, next: Next) -> Response {
    println!("INCOMING REQUEST: {} {}", req.method(), req.uri().path());
    next.run(req).await
}

#[tokio::main]
async fn main() {
    let port_key = "FUNCTIONS_CUSTOMHANDLER_PORT";
    let port: u16 = env::var(port_key)
        .unwrap_or("3000".to_string())
        .parse()
        .expect("Port must be a number");

    // FIX: Ensure this matches '/api/{short_code}'
    let app = Router::new()
        .route("/api/:short_code", get(handlers::handle_redirect))
        .layer(middleware::from_fn(log_request)); // Add logging

    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    println!("Listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
