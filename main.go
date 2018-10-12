package main

import (
	"context"

	"github.com/rexray/gocsi"

	"github.com/jicahoo/csi-unity/provider"
	"github.com/jicahoo/csi-unity/service"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	gocsi.Run(
		context.Background(),
		service.Name,
		"A description of the SP",
		"",
		provider.New())
}
