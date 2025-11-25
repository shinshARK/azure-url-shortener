-- Create Links table
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='Links' and xtype='U')
BEGIN
    CREATE TABLE Links (
        ShortCode NVARCHAR(20) PRIMARY KEY,
        OriginalUrl NVARCHAR(2048) NOT NULL, -- Renamed from LongUrl
        UserID INT NULL,
        CreatedAt DATETIME DEFAULT GETDATE(),
        ExpiresAt DATETIME NULL,
        ClickCount INT DEFAULT 0,
        CustomAlias NVARCHAR(50) NULL,
        IsActive BIT DEFAULT 1 -- Added IsActive
    );
    
    -- Index for UserID for faster lookups
    CREATE INDEX IX_Links_UserID ON Links(UserID);
END
GO
