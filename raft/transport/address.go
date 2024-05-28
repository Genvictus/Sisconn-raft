package transport

import "fmt"

type address struct {
    IP   string
    Port int
}

func NewAddress(ip string, port int) address {
    return address{IP: ip, Port: port}
}

func (a *address) String() string {
    return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

func (a address) Equal(other *address) bool {
    return a.IP == other.IP && a.Port == other.Port
}
