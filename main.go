package main

import (
	"fmt"
	"github.com/Murray-LIANG/gounity"
)

func main() {
    unity, err := gounity.NewUnity("UnityMgmtIP", "username", "password", true)
    if err != nil {
        panic(err)
    }
    fmt.Println("Hello, World")
}
