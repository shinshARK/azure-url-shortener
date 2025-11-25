use anyhow::Result;
use std::env;
use std::time::{SystemTime, UNIX_EPOCH};
use reqwest;
use hmac::{Hmac, Mac};
use sha2::Sha256;
use base64::prelude::*;
use crate::models::AnalyticsEvent;

pub async fn push_to_queue(event: AnalyticsEvent) -> Result<()> {
    let conn_str = env::var("ServiceBusConnectionString")?;
    let (endpoint, key_name, key) = parse_conn_str(&conn_str)?;
    let queue_name = "analytics-queue";
    
    let url = format!("https://{}/{}/messages", endpoint, queue_name);
    let sas_token = generate_sas_token(&endpoint, &queue_name, &key_name, &key)?;

    let client = reqwest::Client::new();
    let res = client.post(&url)
        .header("Authorization", sas_token)
        .header("Content-Type", "application/json")
        .json(&event)
        .send()
        .await?;

    if !res.status().is_success() {
        let error_text = res.text().await?;
        anyhow::bail!("Service Bus Error: {}", error_text);
    }

    // Use eprintln for logging so it shows up in Azure logs
    println!("Analytics sent for {}", event.short_code);
    Ok(())
}

fn generate_sas_token(uri: &str, queue: &str, key_name: &str, key: &str) -> Result<String> {
    let target_uri = urlencoding::encode(&format!("https://{}/{}", uri, queue)).into_owned();
    
    let expiration = SystemTime::now().duration_since(UNIX_EPOCH)?.as_secs() + 3600;
    let to_sign = format!("{}\n{}", target_uri, expiration);
    
    type HmacSha256 = Hmac<Sha256>;
    let mut mac = HmacSha256::new_from_slice(key.as_bytes()).expect("Key invalid");
    mac.update(to_sign.as_bytes());
    let signature = BASE64_STANDARD.encode(mac.finalize().into_bytes());
    
    let signature_encoded = urlencoding::encode(&signature).into_owned();
    
    Ok(format!("SharedAccessSignature sr={}&sig={}&se={}&skn={}", 
        target_uri, signature_encoded, expiration, key_name))
}
fn parse_conn_str(conn_str: &str) -> Result<(String, String, String)> {
    let parts: std::collections::HashMap<_, _> = conn_str
        .split(';')
        .filter_map(|s| {
            let mut split = s.splitn(2, '=');
            Some((split.next()?, split.next()?))
        })
        .collect();

    let endpoint = parts.get("Endpoint")
        .ok_or(anyhow::anyhow!("Missing Endpoint"))?
        .replace("sb://", "")
        .trim_end_matches('/')
        .to_string();
    let key_name = parts.get("SharedAccessKeyName").unwrap_or(&"").to_string();
    let key = parts.get("SharedAccessKey").unwrap_or(&"").to_string();

    Ok((endpoint, key_name, key))
}
