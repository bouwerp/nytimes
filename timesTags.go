package nytimes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const TimesTagsUrl = "https://api.nytimes.com/svc/suggest/v1/timestags"

type TagType string

// descriptor
const Des TagType = "Des"

// geographical location
const Geo TagType = "Geo"

// organisation
const Org TagType = "Org"

// person
const Per TagType = "Per"

// title
const Ttl TagType = "Ttl"

type ListTagsRequest struct {
	Query  string
	Filter TagType
	Max    int64
}

type Tag struct {
	Type  TagType
	Value string
}

type ListTagsResponse struct {
	Tags []Tag
}

func (r ListTagsRequest) validate() error {
	if r.Query == "" {
		return QueryMustBeProvided{}
	}
	return nil
}

func (c *Client) ListTags(request ListTagsRequest) (*ListTagsResponse, error) {
	// validate request
	if err := request.validate(); err != nil {
		return nil, err
	}

	// construct URL
	u, err := url.Parse(TimesTagsUrl)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("query", request.Query)
	if request.Filter != "" {
		query.Add("filter", "("+string(request.Filter)+")")
	}
	if request.Max != 0 {
		query.Add("max", strconv.FormatInt(request.Max, 10))
	}
	query.Add("api-key", c.Key)
	u.RawQuery = query.Encode()

	// execute the request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := (&http.Client{}).Do(req)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	var response []interface{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Fatal(err.Error())
	}

	// decode the response
	tagsList := response[1].([]interface{})
	var tags []Tag
	for _, t := range tagsList {
		reg := regexp.MustCompile("(.*)\\s+\\((Des|Geo|Org|Per|Ttl)\\)")
		m := reg.FindAllStringSubmatch(t.(string), -1)
		if len(m) > 0 && len(m[0]) > 2 {
			tags = append(tags, Tag{
				Type:  TagType(m[0][2]),
				Value: m[0][1],
			})
		} else {
			log.Println("could not parse tag '", t, "'")
		}
	}

	return &ListTagsResponse{Tags: tags}, nil
}

// further tag parser methods

// Person returns a person structure if the tag type is "Per", and an error otherwise.
func (t Tag) Person() (*PersonTag, error) {
	if t.Type != Per {
		return nil, IncorrectTagType{Type: Per}
	}
	names := strings.Split(t.Value, ",")
	return &PersonTag{
		FirstNames: strings.TrimSpace(names[0]),
		LastName:   strings.TrimSpace(names[1]),
	}, nil
}

type PersonTag struct {
	FirstNames string
	LastName   string
}
