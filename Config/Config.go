package Config

import (
	"encoding/json"
	"log"
	"os"
	"telegram_test_bot/Models"
)

func GetConfig() Models.Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}
	decoder := json.NewDecoder(file)
	config := Models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}
	return config
}
