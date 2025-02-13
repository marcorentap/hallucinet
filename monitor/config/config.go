package config

import (
	"os"

	"github.com/marcorentap/hallucinet/types"
)

func getEnvOrDefault(key string, fallback string) string {
	val, err := os.LookupEnv(key)
	if err == false {
		return fallback
	}
	return val
}

// Config is immutable
func NewHallucinetConfig() types.HallucinetConfig {
	return types.HallucinetConfig{
		NetworkName:  getEnvOrDefault("NETWORKNAME", "hallucinet"),
		SqlitePath:   getEnvOrDefault("SQLITE_PATH", "/var/hallucinet/hallucinet.db"),
		DomainSuffix: getEnvOrDefault("DOMAIN_SUFFIX", ".test"),
		HostsPath:    getEnvOrDefault("HOSTS_PATH", "/var/hallucinet/hosts"),
	}
}
