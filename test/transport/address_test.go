package test

import (
	"Sisconn-raft/raft/transport"
	"fmt"
	"testing"
)

func TestAddress(t *testing.T) {
	address := transport.Address {
		IP: "localhost",
		Port: 6969,
	}

	if address.IP != "localhost" {
		t.Error("IP Address not destructed properly, expected", "localhost", "got", address.IP)
	}

	if address.Port != 6969 {
		t.Error("Port Address not destructed properly, expected", 6969, "got", address.Port)
	}
}

func TestNewAddress(t *testing.T) {
	address := transport.NewAddress("localhost", 6969)

	if address.IP != "localhost" {
		t.Error("IP address not destructed properly, expected", "localhost", "got", address.IP)
	}

	if address.Port != 6969 {
		t.Error("Port address not destructed properly, expected", 6969, "got", address.Port)
	}
}

func TestAddress_String(t *testing.T) {
	address := transport.NewAddress("localhost", 6969)

	addressString := address.String()
	if addressString != "localhost:6969" {
		t.Error("String address is incorrect, expected", "localhost:6969", "got", addressString)
	}

	addressStringFmt := fmt.Sprint(&address)
	if addressStringFmt != "localhost:6969" {
		t.Error("fmt string address is incorrect, expected", "localhost:6969", "got", addressStringFmt)
	}
}

func TestAddress_Equal(t *testing.T) {
	address1 := transport.Address {
		IP: "localhost",
		Port: 6969,
	}
	address2 := transport.NewAddress("localhost", 6969)
	address3 := transport.NewAddress("lcoalhost", 6420)

	if address1 != address2 {
		t.Error("Address", &address1, "should be equal to", &address2)
	}

	if !address1.Equal(&address2) {
		t.Error("Address", &address1, "should be equal to", &address2)
	}
	
	if address1 == address3 {
		t.Error("Address", &address1, "should not be equal to", &address3)
	}

	if address1.Equal(&address3) {
		t.Error("Address", &address1, "should not be equal to", &address3)
	}
}
