package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	t.Run("Given env variable returns when populated", func(t *testing.T) {
		f, err := os.CreateTemp(".", ".env.foo")
		if err != nil {
			t.Fail()
		}
		f.WriteString("FOO=foo")
		f.Sync()
		defer os.Remove(f.Name())

		Load(f.Name())

		res := getEnv("FOO", "bar")
		require.Equal(t, res, "foo", "environment variable should match that in .env.test")
	})

	t.Run("Default env value used when env variable unpopulated", func(t *testing.T) {
		res := getEnv("foo", "bar")
		require.Equal(t, res, "bar", "environment variable should not be defined in .env.test")
	})
}

func TestLoad(t *testing.T) {
	t.Run("Config values should use default values if env file empty", func(t *testing.T) {
		f, err := os.CreateTemp(".", ".env.foo")
		if err != nil {
			t.Fail()
		}
		defer os.Remove(f.Name())

		cfg, err := Load(f.Name())
		require.Nil(t, err, "should not get any errors")

		require.Equal(t, cfg.Server.Port, "8080")
		require.Equal(t, cfg.Server.Host, "0.0.0.0")

		require.Equal(t, cfg.Database.Host, "localhost")
		require.Equal(t, cfg.Database.Port, "5432")
		require.Equal(t, cfg.Database.User, "postgres")
		require.Equal(t, cfg.Database.Password, "")
		require.Equal(t, cfg.Database.DBName, "postgres")
		require.Equal(t, cfg.Database.SSLMode, "disable")

		require.Equal(t, cfg.API.BasePath, "/api")

		require.Equal(t, cfg.JWT.AccessTokenSecret, "your-secret-key")
		require.Equal(t, cfg.JWT.RefreshTokenSecret, "your-secret-refresh-key")

		require.Equal(t, cfg.Environment, "development")
	})

	t.Run("Config values should match values in environment file", func(t *testing.T) {
		f, err := os.CreateTemp(".", ".env.foo")
		if err != nil {
			t.Fail()
		}
		f.WriteString("JWT_TOKEN_SECRET=test-jwt-secret\nENV=test")
		f.Sync()
		defer os.Remove(f.Name())

		cfg, err := Load(f.Name())
		require.Nil(t, err, "should not get any errors")

		require.Equal(t, cfg.JWT.AccessTokenSecret, "test-jwt-secret")
		require.Equal(t, cfg.Environment, "test")
	})
}

func TestGetDSN(t *testing.T) {
	t.Run("DSN string should print correctly", func(t *testing.T) {
		var c Config
		c.Database.Host = "host"
		c.Database.Port = "port"
		c.Database.User = "user"
		c.Database.Password = "password"
		c.Database.DBName = "dbname"
		c.Database.SSLMode = "sslmode"

		expected := "host=host port=port user=user password=password dbname=dbname sslmode=sslmode"
		require.Equal(t, c.GetDSN(), expected, "DSN strings should match")
	})
}
