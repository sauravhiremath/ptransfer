package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// BUFFERSIZE is the minimum packet size for file transfers
const BUFFERSIZE = 1024
const numberConnections = 8

// IV is the key for AES pls
var IV = []byte("1234567812345678")

// Build flags
var server, file string

// Global varaibles
var serverAddress, aesKey, fileName string

func main() {
	flag.StringVar(&aesKey, "aesKey", "", "(required) Key or Password to encrypt using AES")
	flag.StringVar(&serverAddress, "server", "", "(run as client) server address to connect to")
	flag.StringVar(&fileName, "file", "", "(run as server) file to serve")
	flag.Parse()
	// Check build flags too, which take precedent
	if server != "" {
		serverAddress = server
	}
	if file != "" {
		fileName = file
	}
	fmt.Println(`
───────────▄███████████████▄───────────────────────────────────────────────────────
────────▄█████████████████████▄────────────────────────────────────────────────────
──────▄███████▀─────────▀████████──────────────────────────────────────────────────
─────███████────────────────██████─────────────────────────────────────────────────
────██████────▄████████▄─────██████────────────────────────────────────────────────
───██████────████████████▄────██████───────────────────────────────────────────────
──██████────███▀───▀██████─────██████▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
─███████───███▌─────███████────██████████████████████████████████░░░░██████████████
████████───███──────███████────████░█████████████████████████████░█████████████████
████████───███─────▐███████────████░███████████████████████████░░░░░███████████████
████████───███▌────▐██████─────██░░░░░██░░░░█░░░░░█░░░░░░█░░░░░██░██░░░░░░░█░█░░░██
████████▄──▀██▌────███████────█████░█████░░██░███░█░░███░█░██████░██░█████░█░░░█░██
█████████▄──█▌─────██████─────█████░█████░███████░█░████░█░░░░░██░██░█░░░░██░██████
█████████████▌─────▀███▀─────██████░░░░██░██░░░░░░█░████░█████░██░██░░██████░██████
█████████████▌─────────────▄███████░░████░██░░░░░░█░████░█░░░░░██░███░░░░░░█░██████
▀████████████▌───▄───────▄███████████▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
─▀███████████▌───████████████████████──────────────────────────────────────────────
──▀██████████▌───███████████████████───────────────────────────────────────────────
───▀█████████▌───██████████████████────────────────────────────────────────────────
────▀████████▌───█████████████████─────────────────────────────────────────────────
─────▀███████▌──▐████████████████──────────────────────────────────────────────────
───────▀██████▄▄███████████████▀───────────────────────────────────────────────────
─────────▀██████████████████▀──────────────────────────────────────────────────────
────────────▀████████████▀─────────────────────────────────────────────────────────
	`)

	if len(fileName) != 0 {
		runServer()
	} else if len(serverAddress) != 0 {
		runClient()
	} else {
		fmt.Println("You must specify either -file (for running as a server) or -server (for running as a client)")
	}
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}