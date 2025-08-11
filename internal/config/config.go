package config

type Config struct {
	Bot struct {
		AgingPeriod int `yaml:aging_period`
	}
}

var AppConfig *Config

func LoadConfig() error {

}
