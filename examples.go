package main

import "fmt"

const username = ""
const password = ""
const baseURL = ""

func search() {
	tp := BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client, err := NewClient(tp.Client(), baseURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	ims, _, err := client.Incident.Search("(IMTicketStatus=\"Open\" or IMTicketStatus=\"Reopened\" or IMTicketStatus=\"Work in progress\")", nil)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	for _, im := range ims {
		fmt.Println(im)
	}
}

func get() {
	tp := BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client, err := NewClient(tp.Client(), baseURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	im, _, err := client.Incident.Get("IM100045")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	fmt.Println(im)
}
