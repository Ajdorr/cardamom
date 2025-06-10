package config

import (
	"cardamom/core/source/ext"
	"embed"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

//go:embed resources/*
var resources embed.FS

var C Config

type Config struct {

	// Server config
	Server struct {
		Env            string `yaml:"env"`
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		Domain         string `yaml:"domain"`
		JwtTokenSecret string `yaml:"jwtTokenSecret"`
		PasswordSalt   string `yaml:"passwordSalt"`

		Log struct {
			Indent string `yaml:"indent"`
		} `yaml:"log"`
	} `yaml:"server"`

	// Events
	Events struct {
		EventFileStreamDirectory string `yaml:"eventFileStreamDirectory"`
	} `yaml:"events"`

	// Database
	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	DB struct {
		// DB_Sqlite string `env:"DB_SQLITE_PATH" envDefault:"/tmp/core.db"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`

	// OAuth2
	OAuth2 struct {
		Github struct {
			ClientId     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		} `yaml:"github"`
		Microsoft struct {
			ClientId     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		} `yaml:"microsoft"`
		Facebook struct {
			ClientId     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		} `yaml:"facebook"`
		Google struct {
			ClientId     string   `yaml:"clientId"`
			ClientSecret string   `yaml:"clientSecret"`
			AuthURI      string   `yaml:"authUri"`
			TokenURI     string   `yaml:"tokenUri"`
			RedirectURIs []string `yaml:"redirectUris"`
		} `yaml:"google"`
	} `yaml:"oauth2"`
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.Username, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Database,
	)
}

// FIXME remove
func IsLocal() bool {
	return C.Server.Env == "local"
}

func LoadConfig(filename string) error {
	data, err := resources.ReadFile("resources/" + filename)
	if err != nil {
		return err
	}

	ext.PanicIfError(yaml.Unmarshal(data, &C))

	return nil
}

func loadEnvConfig() {
	config, ok := os.LookupEnv("CARDAMOM_CONFIG")
	if !ok {
		return
	}

	data, err := os.ReadFile(config)
	ext.PanicIfError(err)
	ext.PanicIfError(yaml.Unmarshal(data, &C))
}

func init() {
	gin.EnableJsonDecoderDisallowUnknownFields()

	// load config
	ext.PanicIfError(LoadConfig("config.yaml"))

	// load override config
	LoadConfig("override.yaml")

	// load config from environment
	loadEnvConfig()
}
