package main

import (
	"fmt"
	"testing"

	"github.com/olebedev/config"
)

func Test_GetConf(t *testing.T) {
	cfg, err := config.ParseYamlFile("./conf.yaml")
	if err != nil {
		panic(err)
	}
	//fmt.Println(host['development']['database']['host'])
	test, _ := cfg.String("production.database.host")
	fmt.Println(test)

}
