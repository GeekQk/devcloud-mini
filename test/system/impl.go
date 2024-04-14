package system

type Object interface {
	GetInfo() string
}

var StrMapObj = make(map[string]Object)

func init() {
	StrMapObj["host"] = &Host{}
}

type Hosts interface {
	CreateHost(host string) error
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
