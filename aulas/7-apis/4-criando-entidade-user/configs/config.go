package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Conf struct {
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

var config *Conf

func init() {
	viper.SetConfigName("app_confi")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
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
}

func GetDBDriver() string {
	return config.DBDriver
}

func GetDBHost() string {
	return config.DBHost
}

func GetDBPort() string {
	return config.DBPort
}

func GetDBUser() string {
	return config.DBUser
}

func GetDBPassword() string {
	return config.DBPassword
}

func GetDBName() string {
	return config.DBName
}

func GetWebServerPort() string {
	return config.WebServerPort
}

func GetJWTSecret() string {
	return config.JWTSecret
}

func GetJWTExpiresIn() int {
	return config.JWTExpiresIn
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return config.TokenAuth
}
