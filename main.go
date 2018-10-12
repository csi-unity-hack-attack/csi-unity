package main

import (
	"fmt"
	"context"

	"github.com/rexray/gocsi"

	"github.com/jicahoo/csi-unity/provider"
	"github.com/jicahoo/csi-unity/service"
	"github.com/Murray-LIANG/gounity"
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
		"A description of the SP",
		"",
		provider.New())
}
