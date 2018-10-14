package main

import (
	"context"
	"fmt"

	"github.com/rexray/gocsi"

	"github.com/jicahoo/csi-unity/provider"
	"github.com/jicahoo/csi-unity/service"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	fmt.Println("Start to initialize csi-unity.")
	gocsi.Run(
		context.Background(),
		service.Name,
		"CSI plugin for Dell EMC Unity storage system.",
		"",
		provider.New())
}
