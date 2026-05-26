package config

import "os"

type Config struct {
	ApiKey      string
	MySQLDSN    string
	StoragePath string
	BaseURL     string
}

func Load() Config {
	return Config{
		ApiKey:      getEnv("BAILIAN_API_KEY", ""),
		MySQLDSN:    getEnv("MYSQL_DSN", "bailian:bailian_pass_2024@tcp(db:3306)/bailian?charset=utf8mb4&parseTime=true"),
		StoragePath: getEnv("STORAGE_PATH", "/data/storage"),
		BaseURL:     "https://dashscope.aliyuncs.com",
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}