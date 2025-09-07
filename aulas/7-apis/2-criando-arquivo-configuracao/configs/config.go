package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth
}

var config *conf

// poderia fazer uma funcao init pra rodar antes do main
func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_confi")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	// AutomaticEnv carrega variáveis de ambiente em detrimento de variáveis de configuração(.env)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// O Unmarshal vai mapear as variáveis do arquivo de configuração para a struct conf
	// O Unmarshal no Go serve para decodificar dados de um formato (como JSON ou YAML) para uma estrutura Go
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	config.TokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)

	return config, err
}
