package system

import "fmt"

func (i *Host) CreateHost(host string) error {
	fmt.Println("create host:", host)
	return nil
}
