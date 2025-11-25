use anyhow::Result;
use std::env;
use tiberius::{Client, Config};
use tokio::net::TcpStream;
use tokio_util::compat::TokioAsyncWriteCompatExt;

// In src/db.rs

pub async fn get_original_url(short_code: &str) -> Result<Option<String>> {
    println!("DEBUG: Querying for ShortCode: '{}'", short_code); // 'Quotes' reveal spaces

    let conn_str = env::var("SqlConnectionString")?;
    let mut config = Config::from_ado_string(&conn_str)?;

    config.encryption(tiberius::EncryptionLevel::Required);
    config.trust_cert(); // Trust Azure's certificate
                         //
    let tcp = TcpStream::connect(config.get_addr()).await?;
    tcp.set_nodelay(true)?;
    
    let mut client = Client::connect(config, tcp.compat_write()).await?;

    // DEBUG: Removed "AND IsActive = 1" to isolate the problem
    let query = "SELECT OriginalUrl, IsActive FROM Links WHERE ShortCode = @P1";
    
    let stream = client.query(query, &[&short_code]).await?;
    let row = stream.into_row().await?;

    if let Some(r) = row {
        let url: &str = r.get(0).unwrap_or("NO_URL");
        let is_active: bool = r.get(1).unwrap_or(false); // Check what the DB actually thinks
        
        println!("DEBUG: Found Row! URL: '{}', IsActive: {}", url, is_active);
        
        // Manual filter in Rust instead of SQL (for debugging)
        if is_active {
            return Ok(Some(url.to_string()));
        } else {
            println!("DEBUG: Row found but IsActive is FALSE");
            return Ok(None);
        }
    }

    println!("DEBUG: No matching ShortCode found in DB.");
    Ok(None)
}
