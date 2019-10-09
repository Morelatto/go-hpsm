package main

import "fmt"
import "github.com/Morelatto/go-hpsm"

func main() {
	const username = ""
	const password = ""
	const baseURL = ""

	tp := hpsm.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client, err := hpsm.NewClient(tp.Client(), baseURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	im, _, err := client.Incident.Get("IM13612124")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	fmt.Println(im)
}
