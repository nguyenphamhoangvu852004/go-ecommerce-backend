package initialize

import (
	"fmt"
	"go-ecommerce-backend-api/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./config")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")

	// // đọc config
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read in YAML file: %w", err))
	}

	fmt.Printf("Raw config: %#v\n", viper.AllSettings())
	// fmt.Println("Server Port: ", viper.GetInt("server.port"))
	// fmt.Println("Database Name: ", viper.GetString("databases.mysql.dbname"))

	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("failed to unmarshal YAML file: %w", err))
	}
}
