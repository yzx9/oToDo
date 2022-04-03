package config

import (
	"github.com/spf13/viper"
)

func SetConfig(config *viper.Viper) {
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
		c := config.Sub("session")
		Session = session{
			AccessTokenExpiresIn:         c.GetInt("access_token_exp"),
			RefreshTokenExpiresInDefault: c.GetInt("refresh_token_exp_default"),
			RefreshTokenExpiresInMax:     c.GetInt("refresh_token_exp_max"),
			RefreshTokenExpiresInOAuth:   c.GetInt("refresh_token_exp_oauth"),
			AccessTokenRefreshThreshold:  c.GetInt("access_token_refresh_threshold"),
		}
	}

	{
		c := config.Sub("secret")
		Secret = secret{
			TokenIssuer:     c.GetString("token_issuer"),
			TokenHmacSecret: []byte(c.GetString("token_hmac_secret")),
			PasswordNonce:   []byte(c.GetString("password_nonce")),
		}
	}

	{
		c := config.Sub("github")
		GitHub = github{
			ClientID:            c.GetString("client_id"),
			ClientSecret:        c.GetString("client_secret"),
			OAuthRedirectURI:    c.GetString("oauth_redirect_uri"),
			OAuthStateExpiresIn: c.GetInt("oauth_state_exp"),
		}
	}
}
