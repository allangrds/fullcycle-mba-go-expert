package configs

type conf struct {
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	WebServerPort string
	JWTSecret     string
	JWTExpiresIn  int
}

var config *conf

// poderia fazer uma funcao init pra rodar antes do main
func LoadConfig(path string) (*conf, error) {}
