package main

import (
	"fmt"
	"net/http"
	"time"
)

var statusData StatusData
var configs Configs

func main() {
	configs.UpdateData()

	statusData = StatusData{
		Status:       "follower",
		LeaderIsLive: false,
	}

	go listen()
	sendMessage()
}

func sendMessage() {
	for {
		statusData.LeaderIsLive = false
		time.Sleep(time.Second)

		if statusData.LeaderIsLive == true {
			continue
		}

		if statusData.Status == "leader" {
			if majorityIsAvailable() {
				leaderMessage()
				continue
			} else {
				statusData.Status = "follower"
			}
		}

		time.Sleep(time.Second * 2)
		if !statusData.LeaderIsLive {
			voting()
		}
	}
}

func leaderMessage() {
	client := http.Client{}

	for _, h := range configs.HostsOtherServices {
		req, _ := http.NewRequest("GET", "http://localhost:"+h+"/mp", nil)
		req.Header.Set("server_status", "leader")
		req.Header.Set("hostAddress", configs.Host)
		go client.Do(req)
	}
}

func majorityIsAvailable() bool {
	countAvailable := 0
	timeout := time.Second / 10
	client := http.Client{
		Timeout: timeout,
	}

	for _, h := range configs.HostsOtherServices {
		req, _ := http.NewRequest("GET", "http://localhost:"+h+"/ping", nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			countAvailable++
		}
	}

	return configs.CountServices/2 < countAvailable
}

func voting() {
	statusData.Status = "candidate"
	countVoices := 0

	timeout := time.Second / 20
	client := http.Client{
		Timeout: timeout,
	}

	for _, h := range configs.HostsOtherServices {
		req, _ := http.NewRequest("GET", "http://localhost:"+h+"/mp", nil)
		req.Header.Set("server_status", "candidate")
		r, err := client.Do(req)
		if err != nil {
			continue
		}
		if r.Header.Get("voice") == "yes" {
			countVoices++
		}
	}

	if configs.CountServices/2 < countVoices && statusData.Status == "candidate" {
		statusData.Status = "leader"
	} else {
		statusData.Status = "follower"
	}
}

func listen() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/mp", messageProcessing)
	err := http.ListenAndServe(":"+configs.Host, nil)
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
}

func messageProcessing(w http.ResponseWriter, r *http.Request) {
	status := r.Header.Get("server_status")

	switch status {
	case "leader":
		if r.Header.Get("hostAddress") == configs.Host {
			fmt.Println("i'm Leader")
		} else {
			statusData.Status = "follower"
			fmt.Println("i'm follower")
		}
		statusData.LeaderIsLive = true
	case "candidate":
		if statusData.Status != "leader" {
			w.Header().Add("voice", "yes")
		}
		statusData.LeaderIsLive = true
	}
}
