package app

import (
	"math/rand"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var C Config

type Config struct {

	// Server config
	Env            string `env:"ENV" envDefault:"local"`
	Host           string `env:"HOST" envDefault:"localhost"`
	Port           string `env:"PORT" envDefault:"3000"`
	Domain         string `env:"DOMAIN"`
	JwtTokenSecret string `env:"JWT_TOKEN_SECRET" envDefault:"CadamomIsGreat!1234"`
	PasswordSalt   string `env:"PASSWORD_SALT" envDefault:"#SecretCardamomSalt-5555"`

	// Events
	EventFileStreamDirectory string `env:"EVENT_FILE_STREAM_PATH" envDefault:""`

	// Database
	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	DB_Host     string `env:"DB_HOST" envDefault:"localhost"`
	DB_Port     string `env:"DB_PORT" envDefault:"5432"`
	DB_Name     string `env:"DB_NAME" envDefault:"core"`
	DB_Username string `env:"DB_USERNAME" envDefault:"gorm"`
	DB_Password string `env:"DB_PASSWORD" envDefault:"Password1234"`

	// Admin user, maybe omit?
	AdminUserEmail    string `env:"ADMIN_USER_EMAIL" envDefault:"admin@cardamom.com"`
	AdminUserPassword string `env:"ADMIN_USER_PASSWORD" envDefault:"a"`

	// OAuth2
	OAuthGithubClientId        string `env:"OAUTH_GITHUB_CLIENT_ID"`
	OAuthGithubClientSecret    string `env:"OAUTH_GITHUB_CLIENT_SECRET"`
	OAuthMicrosoftClientId     string `env:"OAUTH_MICROSOFT_CLIENT_ID"`
	OAuthMicrosoftClientSecret string `env:"OAUTH_MICROSOFT_CLIENT_SECRET"`
	OAuthFacebookClientId      string `env:"OAUTH_FACEBOOK_CLIENT_ID"`
	OAuthFacebookClientSecret  string `env:"OAUTH_FACEBOOK_CLIENT_SECRET"`
	OAuthGoogleJsonFilepath    string `env:"OAUTH_GOOGLE_CREDS_FILEPATH"`
}

func IsProd() bool {
	return C.Env == "prod"
}

func IsLocal() bool {
	return C.Env == "local"
}

func init() {
	godotenv.Load()
	gin.EnableJsonDecoderDisallowUnknownFields()
	rand.Seed(time.Now().Unix())

	// load config
	if err := env.Parse(&C, env.Options{RequiredIfNoDef: true}); err != nil {
		panic(err)
	}
}
