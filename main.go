package main

import (
	"API-Books/initializer"
	"API-Books/routers"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnectToDB()
}

func main() {
	routers.StartServer().Run()
}
