package env

type DBConfig struct {
	ConnString string `env:"DB_CONFIG,required=true"`
}

type StorageConfig struct {
	StorageHost   string `env:"STORAGE_HOST,required=true"`
	StorageKey    string `env:"STORAGE_KEY,required=true"`
	StorageSecret string `env:"STORAGE_SECRET,required=true"`
	StorageSSL    bool   `env:"STORAGE_SSL,default=false"`
}

type JWTSecret struct {
	JWTSecret string `env:"JWT_SECRET,required=true"`
}

type Config struct {
	DBConfig
	StorageConfig
	JWTSecret
}
