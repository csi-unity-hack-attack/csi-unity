package main

import (
	"context"
	"fmt"
	"github.com/Murray-LIANG/gounity"
	"github.com/jicahoo/csi-unity/provider"
	"github.com/jicahoo/csi-unity/service"
	"github.com/rexray/gocsi"
)

func main() {
	unity, err := gounity.NewUnity("10.141.68.198", "admin", "*****!", true)
	if err != nil {
		panic(err)
	}
	pools, err := unity.GetPools()
	fmt.Println("Hello")
	for _, pool := range pools {
		fmt.Println(pool.Name)
	}

	fmt.Println("Hello, World")

	gocsi.Run(
		context.Background(),
		service.Name,
		"A Unity Container Storage Interface (CSI) Plugin",
		usage,
		provider.New())
}

const usage = `
    X_CSI_UNITY_ENDPOINT
	Unity management IP address
 		

    X_CSI_UNITY_USER
	Unity user name

    X_CSI_UNITY_PASSWORD
	Unity password
`
