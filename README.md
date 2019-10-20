# go-hpsm

HP Service Manager Rest API Go wrapper.

## Requirements

* HP SM 9.35+
* Go 1.13

## Examples

### Get

Retrieves a single incident

``` go
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
im, _, err := client.Incident.Get("IM13696192")
if err != nil {
    fmt.Printf("\nerror: %v\n", err)
    return
}
fmt.Println(im)
```

### Search

Query all incidents 

``` go
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
ims, _, err := client.Incident.Search("(IMTicketStatus=\"Open\" or IMTicketStatus=\"Reopened\" or IMTicketStatus=\"Work in progress\")", nil)
if err != nil {
    fmt.Printf("\nerror: %v\n", err)
    return
}
for _, im := range ims {
    fmt.Println(im)
}
```

### TODO

* Close Incident
* Create Incident
* Create Incident With Attachments
* Resolve incident
* Update incident

---

Most of the code adapted from [go-jira](https://github.com/andygrunwald/go-jira).