package web

import (
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/otodo"
)

func SetConfig(config *viper.Viper) {
	{
		c := config.Sub("server")
		otodo.Conf.Server = otodo.ConfigServer{
			ID:               c.GetString("id"),
			FilePathTemplate: c.GetString("file_path_template"),
		}
	}

	{
		c := config.Sub("database")
		otodo.Conf.Database = otodo.ConfigDatabase{
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
		otodo.Conf.Session = otodo.ConfigSession{
			AccessTokenExpiresIn:        c.GetInt("access_token_exp"),
			RefreshTokenExpiresIn:       c.GetInt("refresh_token_exp"),
			AccessTokenRefreshThreshold: c.GetInt("access_token_refresh_threshold"),
		}
	}

	{
		c := config.Sub("secret")
		otodo.Conf.Secret = otodo.ConfigSecret{
			TokenIssuer:     c.GetString("token_issuer"),
			TokenHmacSecret: []byte(c.GetString("token_hmac_secret")),
			PasswordNonce:   []byte(c.GetString("password_nonce")),
		}
	}
}
