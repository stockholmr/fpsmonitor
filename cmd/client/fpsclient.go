package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

type NetworkAdapter struct {
	Name       null.String `json:"name"`
	MacAddress null.String `json:"mac_address"`
	IPAddress  null.String `json:"ip_address"`
}

type Computer struct {
	ComputerName null.String      `json:"name"`
	Username     null.String      `json:"username"`
	Adapters     []NetworkAdapter `json:"adapters"`
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	pcName, _ := os.Hostname()
	data := new(Computer)
	data.ComputerName = null.NewString(pcName, true)
	data.Username = null.NewString(user.Username, true)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, ifa := range ifaces {
		if ifa.HardwareAddr.String() == "" {
			continue
		}

		if strings.Contains(strings.ToLower(ifa.Name), "bluetooth") {
			continue
		}

		if strings.Contains(strings.ToLower(ifa.Name), "vethernet") {
			continue
		}

		adds, _ := ifa.Addrs()
		ips := ""
		for _, a := range adds {
			if strings.Contains(a.String(), "::") {
				continue
			}
			ips += a.String()
		}

		if strings.Contains(ips, "169.254") {
			continue
		}

		data.Adapters = append(data.Adapters, NetworkAdapter{
			Name:       null.NewString(ifa.Name, true),
			MacAddress: null.NewString(ifa.HardwareAddr.String(), true),
			IPAddress:  null.NewString(ips, true),
		})
	}

	jsonStr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	_, err = client.Post("http://127.0.0.1:8080", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

}
