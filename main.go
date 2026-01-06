package main

import (
	"effective-invention/server"
	"fmt"
)

func main() {
	fmt.Println("[ effective-invention ]( launching )")
	server.ServeGin()
}
