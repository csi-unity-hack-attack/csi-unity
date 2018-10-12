package main

import (
	"context"
	"fmt"

	"github.com/rexray/gocsi"

	"github.com/Murray-LIANG/gounity"
	"github.com/jicahoo/csi-unity/provider"
	"github.com/jicahoo/csi-unity/service"
)

// main is ignored when this package is built as a go plug-in.
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
	fmt.Println("Hello World")
	gocsi.Run(
		context.Background(),
		service.Name,
		"CSI plugin for Dell EMC Unity storage system.",
		"",
		provider.New())
}
