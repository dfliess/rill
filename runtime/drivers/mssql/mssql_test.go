package mssql

import (
	"testing"

	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/stretchr/testify/require"
)

func TestResolveDSN_WithOnlyDSN(t *testing.T) {
	c := &ConfigProperties{
		DSN: "sqlserver://user:pass@localhost:1433?database=dbname&encrypt=disable",
	}

	dsn, err := c.ResolveDSN()
	require.NoError(t, err)
	require.Equal(t, c.DSN, dsn)
}

func TestResolveDSN_WithIndividualFields(t *testing.T) {
	c := &ConfigProperties{
		Host:     "db.example.com",
		Port:     1433,
		User:     "admin",
		Password: "secret",
		Database: "analytics",
		Encrypt:  "mandatory",
	}

	dsn, err := c.ResolveDSN()
	require.NoError(t, err)

	expected := "sqlserver://admin:secret@db.example.com:1433?database=analytics&encrypt=mandatory"
	require.Equal(t, expected, dsn)
}

func TestResolveDSN_WithInstance(t *testing.T) {
	c := &ConfigProperties{
		Host:     "db.example.com",
		User:     "admin",
		Password: "secret",
		Database: "analytics",
		Instance: "SQLEXPRESS",
	}

	dsn, err := c.ResolveDSN()
	require.NoError(t, err)

	expected := "sqlserver://admin:secret@db.example.com/SQLEXPRESS?database=analytics"
	require.Equal(t, expected, dsn)
}

func TestResolveDSN_WithTrustServerCertificate(t *testing.T) {
	c := &ConfigProperties{
		Host:                   "localhost",
		Port:                   1433,
		User:                   "sa",
		Password:               "pass",
		Database:               "mydb",
		Encrypt:                "mandatory",
		TrustServerCertificate: true,
	}

	dsn, err := c.ResolveDSN()
	require.NoError(t, err)

	expected := "sqlserver://sa:pass@localhost:1433?database=mydb&encrypt=mandatory&trustservercertificate=true"
	require.Equal(t, expected, dsn)
}

func TestResolveDSN_SpecialPassword(t *testing.T) {
	c := &ConfigProperties{
		User:     "test_user",
		Password: "Aa1 ~`!@#$%^&*()_+-={}[]|\\;:'<>\",./?",
		Host:     "localhost",
		Port:     1433,
		Database: "test",
	}

	dsn, err := c.ResolveDSN()
	require.NoError(t, err)
	expected := "sqlserver://test_user:Aa1%20~%60%21%40%23$%25%5E&%2A%28%29_+-=%7B%7D%5B%5D%7C%5C;%3A%27%3C%3E%22,.%2F%3F@localhost:1433?database=test"
	require.Equal(t, expected, dsn)

	// verify the password round-trips through go-mssqldb's DSN parser
	cfg, err := msdsn.Parse(dsn)
	require.NoError(t, err)
	require.Equal(t, c.User, cfg.User)
	require.Equal(t, c.Password, cfg.Password)
	require.Equal(t, c.Database, cfg.Database)
}

func TestResolveDSN_WithDSNAndIndividualFields_ShouldError(t *testing.T) {
	c := &ConfigProperties{
		DSN:      "sqlserver://some:dsn@localhost?database=db",
		Host:     "localhost",
		Port:     1433,
		User:     "user",
		Database: "testdb",
	}

	_, err := c.ResolveDSN()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid config")
}

func TestResolveDSN_InvalidEncrypt(t *testing.T) {
	c := &ConfigProperties{
		User:     "sa",
		Password: "pass",
		Host:     "localhost",
		Port:     1433,
		Database: "mydb",
		Encrypt:  "invalid-mode",
	}

	_, err := c.ResolveDSN()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported encrypt")
}
