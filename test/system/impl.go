package system

import "fmt"

func init() {
	fmt.Println("init: impl")
}

type Host struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Idc      string `json:"idc"`
	Zone     string `json:"zone"`
	Describe string `json:"describe"`
	Status   string `json:"status"`
}

func (i *Host) GetInfo() string {
	return i.IP
}
