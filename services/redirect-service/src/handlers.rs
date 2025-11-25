use axum::{
    extract::{Path, Request},
    response::{Redirect, IntoResponse, Response},
};
use chrono::Utc;
use crate::db;
use crate::analytics;
use crate::models::AnalyticsEvent;

pub async fn handle_redirect(Path(short_code): Path<String>, req: Request) -> Response {
    // 1. Extract Info
    let user_agent = req.headers()
        .get("user-agent")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("unknown")
        .to_string();
    
    let ip_address = req.headers()
        .get("x-forwarded-for")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("0.0.0.0")
        .split(':').next().unwrap_or("0.0.0.0")
        .to_string();

    // 2. Query DB
    match db::get_original_url(&short_code).await {
        Ok(Some(url)) => {
            // 3. Spawn Async Analytics (Fire & Forget)
            let code_clone = short_code.clone();
            let url_clone = url.clone();
            
            tokio::spawn(async move {
                let event = AnalyticsEvent {
                    short_code: code_clone,
                    original_url: url_clone,
                    timestamp: Utc::now().to_rfc3339(),
                    user_agent,
                    ip_address,
                };
                if let Err(e) = analytics::push_to_queue(event).await {
                    eprintln!("Analytics Fail: {}", e);
                }
            });

            Redirect::temporary(&url).into_response()
        }
        Ok(None) => (axum::http::StatusCode::NOT_FOUND, "Link not found").into_response(),
        Err(e) => {
            eprintln!("DB Error: {}", e);
            (axum::http::StatusCode::INTERNAL_SERVER_ERROR, "Database Error").into_response()
        }
    }
}
