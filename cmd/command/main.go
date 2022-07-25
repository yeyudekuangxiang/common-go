/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"math/rand"
	"mio/internal/app/cmd"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMilli())
	cmd.Execute()
}
