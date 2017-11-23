package main

import (
	"fmt"

	"github.com/BBQgo/iotsrv/bbqapp"
)

func main() {
	fmt.Println("BBQ Application Server.")
	bbqapp.MainLoop()
}
