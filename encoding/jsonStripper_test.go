package encoding

import(
	"testing"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"log"
)

func TestReader_GetIdentifiers(t *testing.T) {
	reader := newTestReader()

	tests := []string{
		"organizations",
		"next_page",
		"count",
		"end_time",
	}
	reader.GetMembers()
	for _, expect := range tests {
		actual := string(reader.NextMember())
		if expect != actual {
			t.Errorf("TestReader_GetIdentifiers failed: expected: %s actual: %s", expect, actual)
		}
		// Skip value associated with this identifier
		reader.NextMember()
	}
}

func TestReader_Unmarshal(t *testing.T) {
	var organizations []testOrganization
	var count int

	reader := newTestReader()
	reader.GetMembers()


	for current := string(reader.NextMember()); len(current) > 0 ; current = string(reader.NextMember()) {
		switch current {
		case "organizations":
			json.Unmarshal(reader.NextMember(), &organizations)
		case "count":
			count, _ = strconv.Atoi(string(reader.NextMember()))
		}
	}

	if result := len(organizations) == count; result != true {
		t.Errorf("TestReader_Unmarshal failed: expected: %d actual: %d\n", count, len(organizations))
	}

	if result := len(organizations) == count -1; result != false {
		t.Errorf("TestReader_Unmarshal failed: expected: %v actual: %v\n", result, false)
	}
}

type testOrganization struct {}

func newTestReader() (*reader){
	dat, err := ioutil.ReadFile("./test-payload.json")
	if err != nil {
		log.Fatalf("Failed to load test data, aborting tests: %s", err)
	}
	return NewJSONStripper(dat)
}