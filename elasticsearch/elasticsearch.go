package elasticsearch

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

func HelloElasticsearch() {
	config := elasticsearch.Config{
		Username: "elastic",
		Password: "123456",
		Addresses: []string{
			"http://localhost:9200",
			"http://localhost:9300",
		},
	}
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	infoData, err := es.Info()
	if err != nil {
		log.Fatalf("error getting response: %s", err)
	}
	defer infoData.Body.Close()
	log.Println(infoData)

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("my_index"),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatal(err)
	}
	es.Index.WithDocumentType("")
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatal(err)
		} else {
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	es.Indices.Create.WithBody(nil)

	log.Println(strings.Repeat("=", 37))

	log.Print("frea speach")
}

/*
// original JSON document
{
	"name":"xianchaoxing",
	"age":10,
	"gender":"male"
}

// converted into a string with its all values that are seperated by a white space
`xianchaoxing 10 male`

// got highlighted
`xianchaoxing 10 <em>male</em>`

// finally, convert back to a JSON document
{
	"name":"xianchaoxing",
	"age":10,
	"gender":"<em>male</em>"
}
*/

type EsDocumentService struct {
}

// convert an object into a string with only values that are
// sperated with an white space.
func toStringWithOnlyValues(data interface{}, keys []string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	strData := string(jsonData)
	for _, v := range keys {
		strData = strings.ReplaceAll(strData, "\""+v+"\""+":", "")
	}
	return strData, nil
}
