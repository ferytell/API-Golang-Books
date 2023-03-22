package main

import "./Router"

func main() {
	var PORT = ":8080"

	Router.StartServer().Run(PORT)
}
