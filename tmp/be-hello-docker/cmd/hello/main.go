package main

import "fmt"

var (
	version = "1.0.1"
)

func main() {
	fmt.Printf("Hello from Docker! This Go application is running inside a container.\nVersion: %s\n", version)
}
