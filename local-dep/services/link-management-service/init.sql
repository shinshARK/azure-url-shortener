-- 1. Create the Database
CREATE DATABASE UrlShortenerDb;
GO

-- 2. Switch to the new Database
USE UrlShortenerDb;
GO

-- 3. Create the Links table
CREATE TABLE Links (
    ShortCode NVARCHAR(20) PRIMARY KEY,
    LongUrl NVARCHAR(2048) NOT NULL,
    ClickCount INT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE()
);
GO
