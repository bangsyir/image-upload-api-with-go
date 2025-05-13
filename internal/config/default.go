package config

type Config struct {
	Port         int
	DatabasePath string
	UploadDir    string
}

func NewConfig() *Config {
	return &Config{
		Port:         8080,
		DatabasePath: "upload.sqlite",
		UploadDir:    "uploads",
	}
}
