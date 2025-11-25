use axum::{
    extract::{Path, Request},
    handler::Handler,
    response::{Redirect, IntoResponse, Response},
    routing::get,
    Router,
};
use std::env;
use std::net::SocketAddr;
use tiberius::{Client, Config};
use tokio::net::TcpStream;
use tokio_util::compat::TokioAsyncWriteCompatExt;
use anyhow::Result;

#[tokio::main]
async fn main() {
    // 1. Get the port Azure assigns to us (Custom Handler magic)
    let port_key = "FUNCTIONS_CUSTOMHANDLER_PORT";
    let port: u16 = env::var(port_key)
        .unwrap_or_else(|_| "3000".to_string())
        .parse()
        .expect("Custom Handler Port must be a number");

    // 2. Define our router
    // The route must match the folder name in step 1 ("redirect")
    let app = Router::new().route("/redirect/{short_code}", get(handle_redirect));

    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    println!("Listening on {}", addr);

    // 3. Start the server
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn handle_redirect(Path(short_code): Path<String>) -> Response {
    match get_original_url(&short_code).await {
        Ok(Some(url)) => {
            // 302 Redirect
            Redirect::temporary(&url).into_response()
        }
        Ok(None) => {
            // 404 Not Found
            (axum::http::StatusCode::NOT_FOUND, "Link not found").into_response()
        }
        Err(e) => {
            eprintln!("DB Error: {}", e);
            (axum::http::StatusCode::INTERNAL_SERVER_ERROR, "Database Error").into_response()
        }
    }
}

async fn get_original_url(short_code: &str) -> Result<Option<String>> {
    // 1. Parse Connection String from Environment
    // Format: "server=tcp:myserver.database.windows.net;database=mydb;..."
    let conn_str = env::var("SqlConnectionString")?;
    let config = Config::from_ado_string(&conn_str)?;

    // 2. Connect (Pure Rust, no ODBC driver needed!)
    let tcp = TcpStream::connect(config.get_addr()).await?;
    tcp.set_nodelay(true)?; // Disable Nagle's algorithm for speed

    let mut client = Client::connect(config, tcp.compat_write()).await?;

    // 3. Query
    // Note: We use @P1 syntax for parameters in Tiberius
    let query = "SELECT OriginalUrl FROM Links WHERE ShortCode = @P1 AND IsActive = 1";
    
    let stream = client.query(query, &[&short_code]).await?;
    let row = stream.into_row().await?;

    if let Some(r) = row {
        let url: &str = r.get(0).unwrap_or_default();
        return Ok(Some(url.to_string()));
    }

    Ok(None)
}
