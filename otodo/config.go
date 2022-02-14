package otodo

var Conf = Config{}

type Config struct {
	Server   ConfigServer
	Database ConfigDatabase
	Session  ConfigSession
	Secret   ConfigSecret
}

type ConfigServer struct {
	ID               string
	FilePathTemplate string // Support :id, :ext, :name, :path, :date
}

type ConfigDatabase struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	Protocol     string
	DatabaseName string
}

type ConfigSession struct {
	AccessTokenExpiresIn        int
	RefreshTokenExpiresIn       int
	AccessTokenRefreshThreshold int
}

type ConfigSecret struct {
	TokenIssuer     string
	TokenHmacSecret []byte
	PasswordNonce   []byte
}
