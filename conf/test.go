package main

import(
    "fmt"
    "github.com/olebedev/config"
)

func main() {
    cfg, err := config.ParseYamlFile("./conf.yaml")
    if err != nil {
        panic(err)
    }
    //fmt.Println(host['development']['database']['host'])
    test, _ := cfg.String("production.database.host")
    fmt.Println(test)

}
