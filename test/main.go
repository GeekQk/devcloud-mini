package main

import (
	"fmt"
	"test/system"
)

func main() {
	a := system.Host{}
	a.CreateHost("127.0.0.1")

	var ser system.Hosts
	ser = system.StrMapObj["host"].(system.Hosts)
	fmt.Printf("%##v\n", ser)
	ser.CreateHost("test")

}
