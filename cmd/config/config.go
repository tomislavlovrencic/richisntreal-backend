package config

var config Cfg

type Cfg struct {
	App   AppConfig
	MySQL MySQL
}

type AppConfig struct {
	Image    string
	ImageTag string
	Name     string
}

type MySQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
