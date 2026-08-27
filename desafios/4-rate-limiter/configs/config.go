// Package configs carrega a configuração do rate limiter a partir de
// variáveis de ambiente, opcionalmente complementadas por um arquivo .env.
package configs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Conf reúne todas as configurações do rate limiter, exigidas pelo desafio
// via variáveis de ambiente.
type Conf struct {
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                int    `mapstructure:"REDIS_DB"`
	RateLimitIPMax         int64  `mapstructure:"RATE_LIMIT_IP_MAX"`
	RateLimitTokenMax      int64  `mapstructure:"RATE_LIMIT_TOKEN_MAX"`
	RateLimitWindowSeconds int    `mapstructure:"RATE_LIMIT_WINDOW_SECONDS"`
	RateLimitBlockSeconds  int    `mapstructure:"RATE_LIMIT_BLOCK_DURATION_SECONDS"`
}

func (c Conf) Window() time.Duration {
	return time.Duration(c.RateLimitWindowSeconds) * time.Second
}

func (c Conf) BlockDuration() time.Duration {
	return time.Duration(c.RateLimitBlockSeconds) * time.Second
}

// Load lê a configuração de variáveis de ambiente, opcionalmente
// complementadas por um arquivo .env dentro de dir. A ausência do arquivo
// não é tratada como erro: dentro do container, as variáveis chegam apenas
// via docker-compose.yaml.
func Load(dir string) (*Conf, error) {
	v := viper.New()

	v.SetDefault("WEB_SERVER_PORT", "8080")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("RATE_LIMIT_IP_MAX", 10)
	v.SetDefault("RATE_LIMIT_TOKEN_MAX", 100)
	v.SetDefault("RATE_LIMIT_WINDOW_SECONDS", 1)
	v.SetDefault("RATE_LIMIT_BLOCK_DURATION_SECONDS", 300)

	// AutomaticEnv faz variáveis de ambiente reais prevalecerem sobre o
	// arquivo .env, exatamente como em aulas/7-apis e aulas/18.
	v.AutomaticEnv()

	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		v.SetConfigFile(envPath)
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var cfg Conf
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
