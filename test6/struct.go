package main

import (
	"encoding/json"
	"io/ioutil"
)

type StatusData struct {
	Status       string
	LeaderIsLive bool
}

type Configs struct {
	CountServices      int
	Host               string
	HostsOtherServices []string
}

func (s *Configs) UpdateData() {
	buff, err := ioutil.ReadFile("test6/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buff, &s)
	if err != nil {
		panic(err)
	}
	configs.CountServices = len(configs.HostsOtherServices)
}
