package main

import (
	"embed"
	"openspec-visualizer/cmd"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	cmd.Run(assets)
}
