use anyhow::Result;
use std::env;
use tiberius::{Client, Config};
use tokio::net::TcpStream;
use tokio_util::compat::TokioAsyncWriteCompatExt;

pub async fn get_original_url(short_code: &str) -> Result<Option<String>> {
    // 1. Config
    let conn_str = env::var("SqlConnectionString")?;
    let config = Config::from_ado_string(&conn_str)?;

    // 2. Connect
    let tcp = TcpStream::connect(config.get_addr()).await?;
    tcp.set_nodelay(true)?;
    
    let mut client = Client::connect(config, tcp.compat_write()).await?;

    // 3. Query
    let query = "SELECT OriginalUrl FROM Links WHERE ShortCode = @P1 AND IsActive = 1";
    let stream = client.query(query, &[&short_code]).await?;
    let row = stream.into_row().await?;

    if let Some(r) = row {
        let url: &str = r.get(0).unwrap_or_default();
        return Ok(Some(url.to_string()));
    }

    Ok(None)
}
