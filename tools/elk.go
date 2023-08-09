package tools

import (
	"log"
	"net/http"
	"strings"

	// "github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

func init() {
	var client *elasticsearch.Client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// log.Println(res)

	// client.Indices.Create("my_index")

	query := `{ query: { match_all: {} } }`
	re, _ := client.Search(
		client.Search.WithIndex("my_index"),
		client.Search.WithBody(strings.NewReader(query)),
	)
	// esapi.Search()
	log.Println(re)

	defer res.Body.Close()
	// log.Println(res)
}

func HttpExample(w http.ResponseWriter, r *http.Request) {
	// ... # Client usage
}
