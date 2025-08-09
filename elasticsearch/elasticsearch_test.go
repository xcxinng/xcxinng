package elasticsearch

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestHelloElasticsearch(t *testing.T) {
	HelloElasticsearch()
}
func parseKeyValueStringToJSON(s string) (interface{}, error) {
	// Split the input string by whitespace.
	values := strings.Split(s, " ")

	// Create a new array to hold the values.
	var arr []interface{}
	for _, value := range values {
		// Parse each value as JSON and append it to the array.
		var v interface{}
		err := json.Unmarshal([]byte(value), &v)
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
	}

	// Marshal the array as JSON and return it.
	return json.Marshal(arr)
}

func Test_jsonToValues(t *testing.T) {
	a := map[string]interface{}{
		"name":   "xianchaoxing",
		"age":    10,
		"gender": "male",
	}
	ret, err := toStringWithOnlyValues(a, []string{"name", "age", "gender"})
	if err != nil {
		log.Fatal(err)
	}
	t.Log(ret)
}

func TestExample(t *testing.T) {
	example()
}
