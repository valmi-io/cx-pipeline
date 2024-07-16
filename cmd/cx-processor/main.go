package main

import (
	"fmt"

	"github.com/spf13/viper"
	env "github.com/valmi-io/cx-pipeline/internal"
)
func main() {
    env.InitConfig()
    fmt.Println("hello world")
    fmt.Println(viper.Get("APP_BACKEND"))
    fmt.Println(viper.Get("KAFKA_BROKER"))
}