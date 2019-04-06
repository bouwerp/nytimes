package nytimes

import (
	"log"
	"testing"
)

func TestClient_ListTags(t *testing.T) {
	client := Client{Key: *apiKey}
	resp, err := client.ListTags(ListTagsRequest{
		Query:  "cork",
		Filter: Geo,
		Max:    10,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(resp)
}
