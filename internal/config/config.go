package config

import (
	"os"
	"time"

	"github.com/cirius-go/generic/slice"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/cirius-go/portfolio-server/pkg/db"
)

var (
	stage              = os.Getenv("STAGE")
	parameterStorePath = os.Getenv("PARAMETER_STORE_PATH")
)

// HTTPServer represents the HTTP server configuration.
type HTTPServer struct {
	Host         string   `envconfig:"HOST"`
	Port         int      `envconfig:"PORT"`
	AllowOrigins []string `envconfig:"ALLOW_ORIGINS"`
}

// Session config.
type Session struct {
	TTL        time.Duration `envconfig:"TTL"`
	Key        []byte        `envconfig:"KEY"`
	RefreshTTL time.Duration `envconfig:"REFRESH_TTL"`
	RefreshKey []byte        `envconfig:"REFRESH_KEY"`
}

// AssetBucket represents the asset bucket configuration.
type AssetBucket struct {
	Name        string `envconfig:"NAME"`
	TempDir     string `envconfig:"TEMP_DIR"`
	TempExpDays int    `envconfig:"TEMP_EXP_DAYS"`
	ActualDir   string `envconfig:"ACTUAL_DIR"`
}

// Config application.
type Config struct {
	HTTPServer   HTTPServer        `envconfig:"HTTP_SERVER"`
	PGDB         db.PostgresConfig `envconfig:"PGDB"`
	CMSSession   Session           `envconfig:"CMS_SESSION"`
	AssetsBucket AssetBucket       `envconfig:"ASSETS_BUCKET"`
}

// C creates a new default config.
func C() *Config {
	return &Config{
		HTTPServer: HTTPServer{
			Host:         "localhost",
			Port:         3000,
			AllowOrigins: []string{"http://localhost:4000"},
		},
		PGDB: db.PostgresConfig{
			Host:     "localhost",
			Port:     5436,
			Username: "dbadmin",
			Password: "dbadmin",
			Database: "portfolio-server",
			Args: map[string]string{
				"sslmode":         "disable",
				"timezone":        "UTC",
				"connect_timeout": "10",
			},
		},
		CMSSession: Session{
			TTL:        12 * time.Hour,
			Key:        []byte("-#XnQ-0(cNKzKm_v2phHb$;o!V4c-5HG"),
			RefreshTTL: 7 * 24 * time.Hour,
			RefreshKey: []byte("WN*@?5{9wltC)?!^/}mVv2UM?KExuBQ6"),
		},
	}
}

// IsLocal indicates if the server is running locally.
func IsLocal() bool {
	return !slice.Includes(stage, "dev", "uat", "prod")
}

// GetStage returns the stage of the server.
func GetStage() string {
	if stage == "" {
		return "local"
	}
	return stage
}

// Load loads the configuration from environment variables and
// parameter store.
func Load(envFiles ...string) (*Config, error) {
	c := C()

	if IsLocal() {
		if err := LoadLocalConfig(c, envFiles...); err != nil {
			return nil, err
		}

		return c, nil
	}

	if parameterStorePath != "" {
		secrets := map[string]string{}
		if err := LoadFromAPS(os.Getenv("PARAMETER_STORE_PATH"), secrets, nil); err != nil {
			return nil, err
		}

		for k, v := range secrets {
			if err := os.Setenv(k, v); err != nil {
				return nil, err
			}
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}

	return c, nil
}

// LoadLocalConfig loads the configuration from environment variables.
func LoadLocalConfig(cfg *Config, customPaths ...string) error {
	if err := godotenv.Load(customPaths...); err != nil {
		return err
	}

	// load from environment variables.
	return envconfig.Process("", cfg)
}
