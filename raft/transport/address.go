package transport

import "fmt"

type Address struct {
    IP   string
    Port int
}

func NewAddress(ip string, port int) Address {
    return Address{IP: ip, Port: port}
}

func (a *Address) String() string {
    return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

func (a *Address) Equal(other *Address) bool {
    return a.IP == other.IP && a.Port == other.Port
}
