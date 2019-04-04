package nytimes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const ArticleSearchUrl = "https://api.nytimes.com/svc/search/v2/articlesearch.json"

type FacetField string

const DayOfWeek FacetField = "day_of_week"
const DocumentType FacetField = "document_type"
const Ingredients FacetField = "ingredients"
const NewsDesk FacetField = "news_desk"
const PubMonth FacetField = "pub_month"
const PubYear FacetField = "pub_year"
const SectionName FacetField = "section_name"
const Source FacetField = "source"
const SubsectionName FacetField = "subsection_name"
const TypeOfMaterial FacetField = "type_of_material"

type Sort string

const Newest Sort = "newest"
const Oldest Sort = "oldest"
const Relevance Sort = "relevance"

type SearchArticlesRequest struct {
	BeginDate   time.Time // 20060102
	EndDate     time.Time // 20060102
	Facet       bool
	FacetFields FacetField
	FacetFilter bool
	FieldList   []string
	// Lucene syntax: http://www.lucenetutorial.com/lucene-query-syntax.html
	FilterQuery string
	// returns 10 results at a time, thus if response hits > 10, there are more pages
	Page int
	// string search query
	Query string
	Sort  Sort
}

type SearchArticlesResponse struct {
	Status    string `json:"status"`
	Copyright string `json:"copyright"`
	Response  struct {
		Docs   []Doc            `json:"docs"`
		Meta   Meta             `json:"meta"`
		Facets map[string]Terms `json:"facets"`
	} `json:"response"`
}

type Meta struct {
	Hits   int `json:"hits"`
	Offset int `json:"offset"`
	Time   int `json:"time"`
}

type Blog struct {
}

type Legacy struct {
}

type Multimedia struct {
	Rank     int         `json:"rank"`
	Subtype  string      `json:"subtype"`
	Caption  string      `json:"caption"`
	Credit   string      `json:"credit"`
	Type     string      `json:"type"`
	URL      string      `json:"url"`
	Height   int         `json:"height"`
	Width    int         `json:"width"`
	Legacy   Legacy      `json:"legacy"`
	SubType  string      `json:"subType"`
	CropName interface{} `json:"crop_name"`
}

type Headline struct {
	Main          string      `json:"main"`
	Kicker        interface{} `json:"kicker"`
	ContentKicker interface{} `json:"content_kicker"`
	PrintHeadline interface{} `json:"print_headline"`
	Name          string      `json:"name"`
	Seo           interface{} `json:"seo"`
	Sub           interface{} `json:"sub"`
}

type Person struct {
	Firstname    string      `json:"firstname"`
	Middlename   string      `json:"middlename"`
	Lastname     string      `json:"lastname"`
	Qualifier    interface{} `json:"qualifier"`
	Title        interface{} `json:"title"`
	Role         string      `json:"role"`
	Organization string      `json:"organization"`
	Rank         int         `json:"rank"`
}

type Byline struct {
	Original     string      `json:"original"`
	Person       []Person    `json:"person"`
	Organization interface{} `json:"organization"`
}

type Doc struct {
	WebURL         string        `json:"web_url"`
	Snippet        string        `json:"snippet"`
	LeadParagraph  string        `json:"lead_paragraph"`
	Blog           Blog          `json:"blog"`
	Source         string        `json:"source"`
	Multimedia     []Multimedia  `json:"multimedia"`
	Headline       Headline      `json:"headline"`
	Keywords       []interface{} `json:"keywords"`
	PubDate        string        `json:"pub_date"`
	DocumentType   string        `json:"document_type"`
	Byline         Byline        `json:"byline"`
	TypeOfMaterial string        `json:"type_of_material"`
	ID             string        `json:"_id"`
	WordCount      int64         `json:"word_count"`
	Score          float64       `json:"score"`
}

type Terms struct {
	Terms []Term `json:"terms"`
}

type Term struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

func (r SearchArticlesRequest) validate() error {
	return nil
}

func (c *Client) SearchArticles(request SearchArticlesRequest) (*SearchArticlesResponse, error) {
	// validate request
	if err := request.validate(); err != nil {
		return nil, err
	}

	// construct URL
	u, err := url.Parse(ArticleSearchUrl)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("begin_date", request.BeginDate.Format("20060102"))
	query.Add("end_date", request.EndDate.Format("20060102"))
	query.Add("facet", strconv.FormatBool(request.Facet))
	// todo investigate possible comma-separated list
	query.Add("facet_fields", string(request.FacetFields))
	query.Add("facet_filter", strconv.FormatBool(request.FacetFilter))
	query.Add("fl", strings.Join(request.FieldList, ","))
	// todo possibly validate correct Lucene syntax
	query.Add("fq", request.FilterQuery)
	query.Add("q", request.Query)
	query.Add("api-key", c.Key)
	u.RawQuery = query.Encode()

	// execute the request
	fmt.Println("url: ", u.String())
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
	var response SearchArticlesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &response, nil
}
