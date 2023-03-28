package libelastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"time"

	es "github.com/olivere/elastic/v7"
)

func onError(err error) {
	fmt.Println(err)
}

// ElasticClient ...
type ElasticClient interface {
	ElasticBasicActions
	ElasticBulkActions
}

// ElasticBasicActions ...
type ElasticBasicActions interface {
	Search(ctx context.Context, indexName string, option SearchOption, MultiSearchResult bool) (result []byte, err error)
	Store(ctx context.Context, name string, doc interface{}, template *DynamicTemplate) (res *es.IndexResponse, err error)
	Remove(ctx context.Context, indexName string, id string) (res *es.DeleteResponse, err error)
	RemoveIndex(ctx context.Context, indexName ...string) (res *es.IndicesDeleteResponse, err error)
	Ping(ctx context.Context, nodeURL string) (*es.PingResult, int, error)
}

// ElasticBulkActions ....
type ElasticBulkActions interface {
	AddBulkProcessor(bulkProcessor BulkProcessor) (err error)
	BulkStore(ctx context.Context, indexName string, processorName string, docs []interface{}, template *DynamicTemplate) (err error)
}

type BulkProcessor struct {
	Name          string
	Workers       int
	BulkActions   int
	BulkSize      int
	FlushInterval time.Duration
}

type Client struct {
	esclient *es.Client
	Config   ClientConfig
}

type ClientConfig struct {
	BulkProcessors map[string]*es.BulkProcessor
}

type SearchOption struct {
	Query es.Query
	Sort  map[string]bool
	From  int
	Size  int
}

// NewClient ..
func NewClient(elasticUrl, elasticUsername, elasticPassword string) (c ElasticClient, err error) {
	esclient, err := es.NewClient(
		es.SetURL(elasticUrl),
		es.SetSniff(false),
		es.SetHealthcheck(false),
		es.SetRetrier(NewElasticRetrier(3*time.Second, onError)),
		es.SetSniff(false),
		es.SetErrorLog(log.New(os.Stderr, "", log.LstdFlags)),
		es.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		es.SetBasicAuth(elasticUsername, elasticPassword),
	)
	if err != nil {
		return
	}

	client := &Client{
		esclient: esclient,
	}

	client.defaultBulkProcessor()
	c = client

	return
}

// Search ...
func (c *Client) Search(ctx context.Context, indexName string, option SearchOption, MultiSearchResult bool) (result []byte, err error) {
	if option.Size == 0 {
		option.Size = 10 // default size to 10
	}

	search := c.esclient.Search().
		Index(indexName).
		Query(option.Query).
		From(option.From).
		Size(option.Size)

	for k, v := range option.Sort {
		search = search.Sort(k, v)
	}

	searchResult, err := search.Do(ctx)
	if err != nil {
		return
	}

	var sources []map[string]interface{}
	source := make(map[string]interface{})
	if searchResult.Hits.TotalHits.Value > 0 {
		for _, hit := range searchResult.Hits.Hits {
			row := make(map[string]interface{})
			err = json.Unmarshal(hit.Source, &row)
			if err != nil {
				return
			}
			if MultiSearchResult {
				sources = append(sources, row)
			} else {
				source = row
			}
		}
	}

	if MultiSearchResult {
		result, err = json.Marshal(sources)
		if err != nil {
			return
		}
	} else {
		result, err = json.Marshal(source)
		if err != nil {
			return
		}
	}

	return
}

func (c *Client) createMappings(ctx context.Context, indexName string, template *DynamicTemplate) (err error) {
	exist, err := c.esclient.IndexExists(indexName).Do(ctx)
	if err != nil {
		return
	}

	if !exist {
		// create mapping
		_, err = c.esclient.CreateIndex(indexName).
			BodyJson(template).
			Index(indexName).
			Do(ctx)
		if err != nil {
			return
		}
	}

	return

}

// Index index(upsert) document to elastic
func (c *Client) Store(ctx context.Context, name string, doc interface{}, template *DynamicTemplate) (*es.IndexResponse, error) {
	err := c.createMappings(ctx, name, template)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	id := getDocumentID(doc)

	res, err := c.esclient.Index().
		Index(name).
		Id(id).
		BodyJson(doc).Do(ctx)

	if urlErr, ok := err.(*url.Error); ok {
		if urlErr.Err == context.Canceled || urlErr.Err == context.DeadlineExceeded {
			// Proceed, but don't mark the node as dead
			return nil, urlErr.Err
		}
	}

	return res, err
}

// Remove ...
func (c *Client) Remove(ctx context.Context, indexName string, id string) (res *es.DeleteResponse, err error) {
	res, err = c.esclient.Delete().
		Index(indexName).
		Id(id).
		Do(ctx)

	return
}

// RemoveIndex ...
func (c *Client) RemoveIndex(ctx context.Context, indexName ...string) (res *es.IndicesDeleteResponse, err error) {
	res, err = c.esclient.DeleteIndex(indexName...).
		Do(ctx)
	return
}

func (c *Client) RemoveByQuerye(ctx context.Context, indexName string, id string) (res *es.DeleteResponse, err error) {
	res, err = c.esclient.Delete().
		Index(indexName).
		Id(id).
		Do(ctx)

	return
}

func findFieldID(fieldName string) bool {
	return fieldName == "ID" || fieldName == "id" || fieldName == "Id"
}

func getDocumentID(doc interface{}) (id string) {
	val := reflect.ValueOf(doc)
	switch val.Kind() {
	case reflect.Struct:
		fieldByname := val.FieldByNameFunc(findFieldID)

		if fieldByname.IsValid() {
			id = fmt.Sprintf("%s", fieldByname)
		}

		return
	case reflect.Map:
		value := val.MapIndex(reflect.ValueOf("id"))
		if value.IsValid() {
			id = fmt.Sprintf("%s", value)
		}
	}

	return
}

// BulkStore bulk index(upsert) document to elastic
func (c *Client) BulkStore(ctx context.Context, indexName string, processorName string, docs []interface{}, template *DynamicTemplate) (err error) {
	processor, err := c.GetBulkProcessor(processorName)
	if err != nil {
		return
	}

	err = c.createMappings(ctx, indexName, template)
	if err != nil {
		return
	}

	for _, doc := range docs {
		id := getDocumentID(doc)
		bulkUpdateReq := es.NewBulkIndexRequest().Type("_doc").Index(indexName).Id(id).Doc(doc)
		processor.Add(bulkUpdateReq)
	}

	return
}

func (c *Client) newBulkProcessor(processor BulkProcessor) (bulkProcessor *es.BulkProcessor, err error) {
	bulkProcessor, err = c.esclient.BulkProcessor().
		Name(processor.Name).
		Workers(processor.Workers).
		BulkActions(processor.BulkActions).     // commit if # requests reach certain of number
		BulkSize(processor.BulkSize).           // commit when document size reach certain size
		FlushInterval(processor.FlushInterval). // commit every interval of time
		Do(context.Background())
	return
}

func (c *Client) defaultBulkProcessor() {
	bulkProcessor := BulkProcessor{
		Name:          "default",
		Workers:       10,
		BulkActions:   1000,            // flush when reach 1000 requests
		BulkSize:      2 << 20,         // flush when reach 2 MB
		FlushInterval: 1 * time.Second, // flush every 1 seconds
	}

	if c.Config.BulkProcessors == nil {
		c.Config.BulkProcessors = make(map[string]*es.BulkProcessor)
	}

	defaultBulkProcessor, _ := c.newBulkProcessor(bulkProcessor)
	c.Config.BulkProcessors[bulkProcessor.Name] = defaultBulkProcessor
}

// AddBulkProcessor ...
func (c *Client) AddBulkProcessor(bulkProcessor BulkProcessor) (err error) {
	processor, err := c.newBulkProcessor(bulkProcessor)
	if err != nil {
		return err
	}

	c.Config.BulkProcessors[bulkProcessor.Name] = processor
	return
}

// GetBulkProcessor ...
func (c *Client) GetBulkProcessor(name string) (processor *es.BulkProcessor, err error) {
	if name == "" {
		name = "default"
	}

	processor, ok := c.Config.BulkProcessors[name]
	if !ok {
		err = fmt.Errorf("bulk processor with name %s not found", name)
		return
	}

	return
}

// Ping ...
func (c *Client) Ping(ctx context.Context, nodeURL string) (*es.PingResult, int, error) {
	return c.esclient.Ping(nodeURL).Do(ctx)
}
