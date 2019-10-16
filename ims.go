package hpsm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

import (
	"github.com/fatih/structs"
	"github.com/google/go-querystring/query"
	"github.com/trivago/tgo/tcontainer"
)

// IM represents a Service Manager incident.
type IM struct {
	AffectedCI      string                `json:"AffectedCI"`
	Area            string                `json:"Area"`
	Assignee        string                `json:"Assignee"`
	AssignmentGroup string                `json:"AssignmentGroup"`
	Category        string                `json:"Category"`
	ClosedBy        string                `json:"ClosedBy"`
	ClosedTime      time.Time             `json:"ClosedTime"`
	ClosureCode     string                `json:"ClosureCode"`
	Company         string                `json:"Company"`
	Description     []string              `json:"Description"`
	Impact          string                `json:"Impact"`
	IncidentID      string                `json:"IncidentID"`
	Location        string                `json:"Location"`
	OpenTime        time.Time             `json:"OpenTime"`
	OpenedBy        string                `json:"OpenedBy"`
	ProblemType     string                `json:"ProblemType"`
	SLAAgreementID  int                   `json:"SLAAgreementID"`
	Service         string                `json:"Service"`
	Solution        []string              `json:"Solution"`
	Status          string                `json:"Status"`
	Subarea         string                `json:"Subarea"`
	TicketOwner     string                `json:"TicketOwner"`
	Title           string                `json:"Title"`
	UpdatedBy       string                `json:"UpdatedBy"`
	UpdatedTime     time.Time             `json:"UpdatedTime"`
	Urgency         string                `json:"Urgency"`
	OtherFields     tcontainer.MarshalMap // Rest of the fields should go here.
}

// IMService handles IMs for the API.
type IMService struct {
	client *Client
}

type getResult struct {
	IM         IM       `json:"Incident"`
	Messages   []string `json:"Messages"`
	ReturnCode int      `json:"ReturnCode"`
}

// Get returns a full representation of the incident for the given id.
func (s *IMService) Get(incidentID string) (*IM, *http.Response, error) {
	apiEndpoint := fmt.Sprintf("/SM/9/rest/incidents/%s", incidentID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(getResult)
	resp, err := s.client.Do(req, v)
	if err != nil {
		jerr := NewSMError(resp, err)
		return nil, resp, jerr
	}

	return &v.IM, resp, nil
}

// SearchOptions specifies the optional parameters to methods that support pagination.
type SearchOptions struct {
	// Sort: The sort field and order to be returned. {primaryField}:{ascending|descending}[,{secondaryField}]
	Sort string `url:"sort,omitempty"`
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"start,omitempty"`
	// Indicates the number of collection members to be included in the response. Min value: 1.
	Count int `url:"count,omitempty"`
	// Represents a collection. Values: summary, condense (default), expand.
	View string `url:"view,omitempty"`
}

// searchResult is only a small wrapper around the Search (with query) method
// to be able to parse the results
type searchResult struct {
	Total    int      `json:"@totalcount"`
	StartAt  int      `json:"@start"`
	Count    int      `json:"@count"`
	Messages []string `json:"@messages"`
	Results  []struct {
		Incident IM
	} `json:"content"`
	ReturnCode int `json:"ReturnCode"`
}

// Search will search for incidents according to the query
func (s *IMService) Search(q string, options *SearchOptions) ([]IM, *http.Response, error) {
	var u string
	u = fmt.Sprintf("/SM/9/rest/incidents/?query=%s", url.QueryEscape(q))
	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		u += fmt.Sprintf("&%s", q.Encode())
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return []IM{}, nil, err
	}

	v := new(searchResult)
	resp, err := s.client.Do(req, v)
	if err != nil {
		err = NewSMError(resp, err)
	}

	var res []IM
	for _, im := range v.Results {
		res = append(res, im.Incident)
	}
	return res, resp, err
}

// MarshalJSON is a custom JSON marshal function for the IM struct.
// It handles custom fields and maps those from / to "OtherFields" key.
func (i *IM) MarshalJSON() ([]byte, error) {
	m := structs.Map(i)
	unknowns, okay := m["OtherFields"]
	if okay {
		// if other fields are present, shift all key value from unknown to a level up
		for key, value := range unknowns.(tcontainer.MarshalMap) {
			m[key] = value
		}
		delete(m, "OtherFields")
	}
	return json.Marshal(m)
}

// UnmarshalJSON is a custom JSON marshal function for the IM struct.
// It handles custom fields and maps those from / to "OtherFields" key.
func (i *IM) UnmarshalJSON(data []byte) error {

	// Do the normal unmarshalling first
	// Details for this way: http://choly.ca/post/go-json-marshalling/
	type Alias IM
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	totalMap := tcontainer.NewMarshalMap()
	err := json.Unmarshal(data, &totalMap)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagDetail := field.Tag.Get("json")
		if tagDetail == "" {
			// ignore if there are no tags
			continue
		}
		options := strings.Split(tagDetail, ",")

		if len(options) == 0 {
			return fmt.Errorf("no tags options found for %s", field.Name)
		}
		// the first one is the json tag
		key := options[0]
		if _, okay := totalMap.Value(key); okay {
			delete(totalMap, key)
		}

	}
	i = (*IM)(aux.Alias)
	// all the tags found in the struct were removed. Whatever is left are unknowns to struct
	i.OtherFields = totalMap
	return nil
}
