package main

import "fmt"

type Host struct {
	IP       string `json:"ip"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Idc      string `json:"idc"`
	Zone     string `json:"zone"`
	Describe string `json:"describe"`
	Status   string `json:"status"`
}

func main() {
	fmt.Println("hello")
}
