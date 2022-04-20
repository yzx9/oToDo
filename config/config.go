package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/driven/github"
)

var IdentityDomain = identity.Config{}
var GitHubAdapter = github.Config{}

var configManager *viper.Viper

func Load(dir string) error {
	configManager = viper.New()
	configManager.SetConfigType("yaml")
	configManager.AddConfigPath(dir)

	configManager.SetConfigName("config.yaml")
	if err := configManager.ReadInConfig(); err != nil {
		return fmt.Errorf("fails to load config.yaml: %w", err)
	}

	configManager.SetConfigName("secret.yaml")
	if err := configManager.MergeInConfig(); err != nil {
		return fmt.Errorf("fails to load secret.yaml: %w", err)
	}

	setConfig(configManager)

	return nil
}

// TODO[bug]: how to notify config change? eg: github adapter
func LoadAndWatch(dir string) (<-chan time.Time, error) {
	if err := Load(dir); err != nil {
		return nil, err
	}

	changed := make(chan time.Time, 1)
	configManager.OnConfigChange(func(e fsnotify.Event) {
		setConfig(configManager)

		select {
		case changed <- time.Now():
		default:
		}
	})

	configManager.WatchConfig()

	return changed, nil
}

func setConfig(config *viper.Viper) {
	{
		c := config.Sub("server")
		Server = server{
			ID:                       c.GetString("id"),
			Port:                     c.GetInt("port"),
			Host:                     c.GetString("host"),
			AccessControlAllowOrigin: c.GetString("access_control_allow_origin"),
			FilePathTemplate:         c.GetString("file_path_template"),
		}
	}

	{
		c := config.Sub("database")
		Database = database{
			Host:         c.GetString("host"),
			Port:         c.GetInt("port"),
			UserName:     c.GetString("username"),
			Password:     c.GetString("password"),
			Protocol:     c.GetString("protocol"),
			DatabaseName: c.GetString("dbname"),
		}
	}

	{
		c := config.Sub("identity")
		IdentityDomain = identity.Config{
			AccessTokenExpiresIn:         c.GetInt("access_token_exp"),
			AccessTokenRefreshThreshold:  c.GetInt("access_token_refresh_threshold"),
			RefreshTokenExpiresInDefault: c.GetInt("refresh_token_exp_default"),
			RefreshTokenExpiresInMax:     c.GetInt("refresh_token_exp_max"),
			RefreshTokenExpiresInOAuth:   c.GetInt("refresh_token_exp_oauth"),
			OAuthStateExpiresIn:          c.GetInt("oauth_state_exp"),
			TokenIssuer:                  c.GetString("token_issuer"),

			TokenHmacSecret: []byte(c.GetString("token_hmac_secret")),
			PasswordNonce:   []byte(c.GetString("password_nonce")),
		}
	}

	{
		// TODO: remove, used in file domain
		c := config.Sub("identity")
		Secret = secret{
			TokenIssuer:     c.GetString("token_issuer"),
			TokenHmacSecret: []byte(c.GetString("token_hmac_secret")),
			PasswordNonce:   []byte(c.GetString("password_nonce")),
		}
	}

	{
		c := config.Sub("github")
		GitHubAdapter = github.Config{
			ClientID:         c.GetString("client_id"),
			ClientSecret:     c.GetString("client_secret"),
			OAuthRedirectURI: c.GetString("oauth_redirect_uri"),
		}
	}
}
