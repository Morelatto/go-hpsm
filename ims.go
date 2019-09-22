package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)
import "github.com/google/go-querystring/query"

// IM represents a Service Manager incident.
type IM struct {
	AffectedCI      string    `json:"AffectedCI"`
	Area            string    `json:"Area"`
	Assignee        string    `json:"Assignee"`
	AssignmentGroup string    `json:"AssignmentGroup"`
	Category        string    `json:"Category"`
	ClosedBy        string    `json:"ClosedBy"`
	ClosedTime      time.Time `json:"ClosedTime"`
	ClosureCode     string    `json:"ClosureCode"`
	Company         string    `json:"Company"`
	Description     []string  `json:"Description"`
	Impact          string    `json:"Impact"`
	IncidentID      string    `json:"IncidentID"`
	Location        string    `json:"Location"`
	OpenTime        time.Time `json:"OpenTime"`
	OpenedBy        string    `json:"OpenedBy"`
	ProblemType     string    `json:"ProblemType"`
	SLAAgreementID  int       `json:"SLAAgreementID"`
	Service         string    `json:"Service"`
	Solution        []string  `json:"Solution"`
	Status          string    `json:"Status"`
	Subarea         string    `json:"Subarea"`
	TicketOwner     string    `json:"TicketOwner"`
	Title           string    `json:"Title"`
	UpdatedBy       string    `json:"UpdatedBy"`
	UpdatedTime     time.Time `json:"UpdatedTime"`
	Urgency         string    `json:"Urgency"`
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
