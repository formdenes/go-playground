package config

type StorageConfig struct {
	StoragePictureBucket                 string `yaml:"storagePictureBucket" env:"STORAGE_PICTURE_BUCKET" env-required:"true"`
	StorageFileBucket                    string `yaml:"storageFileBucket" env:"STORAGE_FILE_BUCKET" env-required:"true"`
	StoragePictureCleanupIntervalMinutes int    `yaml:"storagePictureCleanupIntervalMinutes" env:"STORAGE_PICTURE_CLEANUP_INTERVAL_MINUTES" env-default:"10"`
	StorageFileCleanupIntervalMinutes    int    `yaml:"storageFileCleanupIntervalMinutes" env:"STORAGE_FILE_CLEANUP_INTERVAL_MINUTES" env-default:"10"`
	AccessKey                            string `yaml:"accessKey" env:"STORAGE_SERVICE_CACCESS_KEY" env-required:"true"`
	SecretKey                            string `yaml:"secretKey" env:"STORAGE_SERVICE_SECRET_KEY" env-required:"true" mask:"fixed"`
	Host                                 string `yaml:"host" env:"STORAGE_SERVICE_HOST" env-default:"localhost"`
	Port                                 string `yaml:"port" env:"STORAGE_SERVICE_PORT" env-default:"9000"`
	Type                                 string `yaml:"type" env:"STORAGE_TYPE" env-default:"minio"`
}

type EmailConfig struct {
	Host     string `yaml:"host" env:"EMAIL_HOST" env-default:"smtp.gmail.com"`
	Port     int    `yaml:"port" env:"EMAIL_PORT" env-default:"587"`
	User     string `yaml:"user" env:"EMAIL_USER" env-required:"true"`
	Password string `yaml:"password" env:"EMAIL_PASSWORD" env-required:"true" mask:"fixed"`
}

type Config struct {
	Server struct {
		Port int    `yaml:"port" env:"SERVER_PORT" env-default:"4000"`
		Host string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	}
	Database struct {
		Port     int    `yaml:"port" env:"DATABASE_PORT" env-default:"5432"`
		Host     string `yaml:"host" env:"DATABASE_HOST" env-default:"localhost"`
		User     string `yaml:"user" env:"DATABASE_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"DATABASE_PASSWORD" env-default:"password" mask:"fixed"`
	}
	Auth struct {
		AuthTokenSecret               string `yaml:"authTokenSecret" env:"AUTH_TOKEN_SECRET" env-default:"secret" mask:"fixed" env-required:"true"`
		AuthTokenIssuer               string `yaml:"authTokenIssuer" env:"AUTH_TOKEN_ISSUER" env-default:"thk"`
		AuthTokenExpirationMinutes    int    `yaml:"authTokenExpirationMinutes" env:"AUTH_TOKEN_EXPIRATION_MINUTES" env-default:"60"`
		RefreshTokenSecret            string `yaml:"refreshTokenSecret" env:"REFRESH_TOKEN_SECRET" env-default:"secret" mask:"fixed"`
		RefreshTokenIssuer            string `yaml:"refreshTokenIssuer" env:"REFRESH_TOKEN_ISSUER" env-default:"thk"`
		RefreshTokenExpirationMinutes int    `yaml:"refreshTokenExpirationMinutes" env:"REFRESH_TOKEN_EXPIRATION_MINUTES" env-default:"60"`
		Google                        struct {
			ClientSecret string `yaml:"clientSecret" env:"AUTH_GOOGLE_CLIENT_SECRET" mask:"fixed" env-required:"true"`
			ClientId     string `yaml:"clientId" env:"AUTH_GOOGLE_CLIENT_ID" env-required:"true"`
		}
	}
	Email   EmailConfig
	Storage StorageConfig
}
