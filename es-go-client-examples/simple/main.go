package main

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

type TestData struct {
	Foo  string    `json:"foo"`
	Time time.Time `json:"time"`
}

func main() {
	t := time.Now()
	index := "es-test-index_" + t.Format("2006-01-02")

	host := TestData{
		Foo:  "Bar",
		Time: t,
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		log.Error(err)
	}

	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Error(err)
	}

	if !exists {
		_, err := client.CreateIndex(index).Do(context.Background())
		if err != nil {
			log.Error(err)
		}
	}

	h, err := client.Index().Index(index).Type("test-data").BodyJson(host).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexed test data with ID %s to index %s, type %s\n", h.Id, h.Index, h.Type)
}
