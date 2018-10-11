package main

import (
	"fmt"
	"github.com/Murray-LIANG/gounity"
)

func main() {
	unity, err := gounity.NewUnity("10.141.68.198", "admin", "Password123!", true)
	if err != nil {
		panic(err)
	}
	pools, err := unity.GetPools()
	fmt.Println("Hello")
	for _, pool := range pools {
		fmt.Println(pool.Name)
	}

	fmt.Println("Hello, World")
}
