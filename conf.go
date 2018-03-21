package main

import (
	"fmt"

	"github.com/olebedev/config"
)

func main() {

}

func test() {
	cfg, err := config.ParseYamlFile("./conf/conf.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	name1, err1 := cfg.String("dev.kafka.0.host")

}
