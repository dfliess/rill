CREATE DATABASE test;
GO
USE test;
GO

CREATE TABLE all_datatypes (
    tinyint_col TINYINT,
    smallint_col SMALLINT,
    int_col INT,
    bigint_col BIGINT,
    bit_col BIT,
    float_col FLOAT,
    real_col REAL,
    decimal_col DECIMAL(10,2),
    numeric_col NUMERIC(12,4),
    money_col MONEY,
    smallmoney_col SMALLMONEY,

    char_col CHAR(10),
    varchar_col VARCHAR(255),
    varcharmax_col VARCHAR(MAX),
    nchar_col NCHAR(10),
    nvarchar_col NVARCHAR(255),
    nvarcharmax_col NVARCHAR(MAX),
    text_col TEXT,
    ntext_col NTEXT,

    binary_col BINARY(10),
    varbinary_col VARBINARY(255),

    date_col DATE,
    time_col TIME,
    datetime_col DATETIME,
    datetime2_col DATETIME2,
    smalldatetime_col SMALLDATETIME,
    datetimeoffset_col DATETIMEOFFSET,

    uniqueidentifier_col UNIQUEIDENTIFIER,
    xml_col XML
);
GO

INSERT INTO all_datatypes (
    tinyint_col, smallint_col, int_col, bigint_col, bit_col,
    float_col, real_col, decimal_col, numeric_col, money_col, smallmoney_col,
    char_col, varchar_col, varcharmax_col, nchar_col, nvarchar_col, nvarcharmax_col, text_col, ntext_col,
    binary_col, varbinary_col,
    date_col, time_col, datetime_col, datetime2_col, smalldatetime_col, datetimeoffset_col,
    uniqueidentifier_col, xml_col
) VALUES
(
    -- Row 1: All NULL values
    NULL, NULL, NULL, NULL, NULL,
    NULL, NULL, NULL, NULL, NULL, NULL,
    NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL,
    NULL, NULL,
    NULL, NULL, NULL, NULL, NULL, NULL,
    NULL, NULL
),
(
    -- Row 2: Non NULL values
    255, 32767, 2147483647, 999999999999999999, 1,
    2.2, 1.5, 3.3, 4.4, 5.6, 7.8,
    'C', 'VarChar', 'Long varchar content', N'N', N'NVarChar', N'Long nvarchar content', 'Text', N'NText',
    CONVERT(BINARY(10), 'Binary'), CONVERT(VARBINARY(255), 'VarBinary'),
    '2024-02-14', '12:34:56', '2025-02-14 12:34:56', '2025-02-14 12:34:56.1234567', '2025-02-14 12:34:00', '2025-02-14 12:34:56 +02:00',
    'AABBCCDD-EEFF-1122-3344-556677889900', '<root><key>value</key></root>'
),
(
    -- Row 3: All Zero values
    0, 0, 0, 0, 0,
    0.0, 0.0, 0.00, 0.0000, 0.0000, 0.0000,
    '', '', '', N'', N'', N'', '', N'',
    CONVERT(BINARY(10), ''), CONVERT(VARBINARY(255), ''),
    '1900-01-01', '00:00:00', '1900-01-01 00:00:00', '1900-01-01 00:00:00', '1900-01-01 00:00:00', '1900-01-01 00:00:00 +00:00',
    '00000000-0000-0000-0000-000000000000', ''
);
GO
