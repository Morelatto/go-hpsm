package main

import "fmt"
import "github.com/Morelatto/go-hpsm"

func main() {
	const username = ""
	const password = ""
	const baseURL = ""

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
