package nytimes

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var apiKey *string

func init() {
	apiKey = flag.String("api-key", "", "NY times API key")
	flag.Parse()
}

func TestClient_SearchArticles(t *testing.T) {
	client := Client{Key: *apiKey}
	beginDate, err := time.Parse("20060102", "20181201")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	endDate, err := time.Parse("20060102", "20190101")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println("executing service")
	resp, err := client.SearchArticles(SearchArticlesRequest{
		BeginDate:   beginDate,
		EndDate:     endDate,
		FieldList:   []string{"headline"},
		Facet:       true,
		FacetFields: []FacetField{DayOfWeek, Source},
		FacetFilter: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(resp)
}
