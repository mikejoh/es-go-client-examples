package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic"
)

const indexMapping = `{
	"mappings" : {
		"ping" : {
			"properties" : {
				"ip" : { "type" : "keyword" },
				"port" : { "type" : "keyword" },
				"status" : { "type" : "keyword" },
				"time" : { "type" : "date" }
			}
		}
	}
}`

type Host struct {
	IP     string    `json:"ip"`
	Port   string    `json:"port"`
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func (h *Host) Ping() (string, error) {
	t := time.Now()
	h.Time = t
	port, _ := strconv.Atoi(h.Port)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.IP, port), 1*time.Second)
	if err != nil {
		return "DOWN", err
	}
	conn.Close()
	return "UP", nil
}

func parseCSV(file string) ([]*Host, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the CSV file
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var hosts []*Host
	// Loop through the CSV
	for _, line := range lines {
		// Convert each row into a struct and append to array of Hosts
		hosts = append(hosts, &Host{
			IP:   line[0],
			Port: line[1],
		})
	}

	return hosts, nil
}

func sendToES(c *elastic.Client, host *Host, index string) error {
	_, err := c.Index().Index(index).Type("ping").BodyJson(&host).Do(context.Background())

	if err != nil {
		return err
	}
	
	return nil
}

func initES(index string) *elastic.Client {
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err := client.CreateIndex(index).Body(indexMapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	return client
}

func main() {
	t := time.Now()
	index := "es-ping_" + t.Format("2006-01-02")

	c := initES(index)

	hosts, _ := parseCSV("hosts.csv")

	for _, host := range hosts {
		status, err := host.Ping()
		if err != nil {
			log.Println(err)
		}
		host.Status = status
		sendToES(c, host, index)
	}
}
