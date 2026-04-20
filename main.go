package main

import (
	"flag"
	"log"

	"github.com/bikash-das/tiny-kv/server"
)

// Global config tied to flags
var Config struct {
	Host string
	Port int
}

func setupFlags() {
	// flag.StringVar and IntVar bind the command line input directly to our struct
	flag.StringVar(&Config.Host, "host", "0.0.0.0", "host for the redis server")
	flag.IntVar(&Config.Port, "port", 7379, "port for the redis server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Initializing tiny-kv...")

	// Pass the flagged values into the server package
	server.RunSyncTCPServer(Config.Host, Config.Port)
}
