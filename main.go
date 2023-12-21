package main

import (
	"log"

	"github.com/ivankuchin/connectopia.org/internal/pkg/server"
)

func SetLogFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	SetLogFlags()

	server.Run()
}
