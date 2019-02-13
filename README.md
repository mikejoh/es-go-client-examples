# Elasticsearch Go-client examples

_I created this repository as a note-to-self_

1. Within the cloned repository run `docker-compose up -d`
2. Run `go run main.go` in any of the sub-folders of the examples, beaware that you need to download the Go-client used in these example and as well as the `logrus` package.
3. To install the missing packages run:
```
go get github.com/sirupsen/logrus
go get github.com/olivere/elastic
```

The rest of the folders within this repository (`kibana` and `es`) are used to seperate configuration files and data directories created during runtime (of e.g. the Elasticsearch data directory).
