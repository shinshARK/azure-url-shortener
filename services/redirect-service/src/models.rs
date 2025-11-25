use serde::Serialize;

#[derive(Serialize)]
pub struct AnalyticsEvent {
    pub short_code: String,
    pub original_url: String,
    pub timestamp: String,
    pub user_agent: String,
    pub ip_address: String,
}
