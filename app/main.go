package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	log "github.com/sirupsen/logrus"
	"mongologs/model"
	"mongologs/parser"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	elasticsearchIndex := os.Getenv("ELASTICSEARCH_INDEX")
	logFile := os.Getenv("LOG_FILE")

	log.Infoln("ELASTICSEARCH_INDEX: ", elasticsearchIndex)
	log.Infoln("LOG_FILE:            ", logFile)
	log.Infoln("ELASTICSEARCH_URL:   ", os.Getenv("ELASTICSEARCH_URL"))

	err := parseFile(logFile, elasticsearchIndex)
	if err != nil {
		log.Error(err)
	}
}

func parseFile(filename string, elasticsearchIndex string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Creates ES client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var i = 0
	var mongologsBulk []*model.MongoLog

	// Stream file and read each line
	for scanner.Scan() {

		// Parses Mongo log
		mongoLog, err := parser.ParseLog(scanner.Text(), i)

		if err != nil {
			log.Warn(err)
			continue
		}

		// Prepare list of logs for bulk operation
		mongologsBulk = append(mongologsBulk, mongoLog)

		// Once 1000 logs reached, store them in ES
		if len(mongologsBulk)%1000 == 0 {
			Bulk(es, elasticsearchIndex, mongologsBulk)
			mongologsBulk = []*model.MongoLog{}
		}

		i = i + 1
	}

	return nil
}

func Bulk(es *elasticsearch.Client, elasticsearchIndex string, mongoLogs []*model.MongoLog) {
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         elasticsearchIndex,
		Client:        es,
		NumWorkers:    2,
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		log.Infoln(3)
		log.Fatalf("Error creating the indexer: %s", err)
	}

	for _, mongoLog := range mongoLogs {
		// Prepare the data payload: encode article to JSON
		//
		mongoLogBytes, _ := json.Marshal(mongoLog)
		mongoLogString := string(mongoLogBytes)
		if err != nil {
			log.Fatalf("Cannot encode mongoLog %d: %s", mongoLog.Id, err)
		}

		err = bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(mongoLog.Id),

				// Body is an `io.Reader` with the payload
				Body: strings.NewReader(mongoLogString),

				//OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					log.Infoln("success")
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					log.Infoln("error")
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	// Check any failures
	bulkIndexerStats := bulkIndexer.Stats()
	if bulkIndexerStats.NumFailed > 0 {
		log.Warnf(
			"Indexed [%d] documents with [%d] errors",
			bulkIndexerStats.NumFlushed,
			bulkIndexerStats.NumFailed,
		)
	}
}
