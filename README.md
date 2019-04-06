# nytimes
Golang implementation of the NY Times API.

# Scope

The state of the various NY times APIs that are available are:
* Article Search: _Implemented_
* Times Tags: _Implemented_

# Usage

```$golang
client := Client{
    Key: "<API_KEY>",
}

beginDate, _ := time.Parse("20060102", "20181201")

endDate, _ := time.Parse("20060102", "20190101")

resp, _ := client.SearchArticles(SearchArticlesRequest{
    BeginDate:   beginDate,
    EndDate:     endDate,
    FieldList:   []string{"headline"},
    Facet:       true,
    FacetFields: []FacetField{DayOfWeek, Source},
    FacetFilter: true,
})

```
