/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/zimlewis/tomato/cmd"
	"github.com/zimlewis/tomato/storage"
)

func main() {

	storage.Initialize()
	cmd.Execute()
}
