package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(args[0], "[new|auth|hi]", "args...")
		os.Exit(1)
	}
	switch args[1] {
	case "new":
		newM()
	case "auth":
		if len(args) < 3 {
			fmt.Println("must provide macaroon")
		}
		authM(args[2])
	case "hi":
		if len(args) < 3 {
			fmt.Println("must provide macaroon")
		}
		hi(args[2])
	default:
		fmt.Println(args[1], "invalid command")
		os.Exit(1)
	}
}

func newM() {
	resp, err := http.Get("http://localhost:9000/new")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}

func authM(serialMacaroon string) {
}

func hi(serialMacaroon string) {
	req, err := http.NewRequest("GET", "http://localhost:9999/hi", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Macaroon", serialMacaroon)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
