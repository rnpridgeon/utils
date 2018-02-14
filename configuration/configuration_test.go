package configuration

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"io/ioutil"
)

var (
	conf *testConfig
)

type testConfig struct {
	StringTest string   `json:"StringTest"`
	IntTest    int      `json:"IntTest"`
	Nested     nestable `json:"Nested"`
}

type nestable struct {
	TestNest string
}

func TestStringArg(t *testing.T) {
	providers := []struct {
		dat       string
		configure func(string) func(config interface{}) error
	}{
		{"test-config.json", FromFile},
		{"http://localhost:8080/test-config.json", FromHTTP},
	}

	var tests = []struct {
		eval       func(interface{}) (want interface{}, result bool)
		input      interface{}
		boolResult bool
	}{
		{testInt, 1, true},
		{testInt, -1, false},
		{testString, "pass", true},
		{testString, "fail", false},
		{testNested, "pass", true},
		{testNested, "fail", false},
	}

	var err error
	for _, provider := range providers {
		if newConf(provider.configure(provider.dat)) != nil {
			t.Fatal("FAILED: %s, %s ", t.Name(), err)
		}
		for _, test := range tests {
			if want, result := test.eval(test.input); result != test.boolResult {
				t.Errorf("FAILED: TestFromFile(conf.json) wanted %v got %v", test.input, want)
			}
		}
	}
}

func TestByteArg(t *testing.T) {
	providers := []struct {
		dat       []byte
		configure func([]byte) func(config interface{}) error
	}{
		{readBytes(), FromBytes},
	}

	var tests = []struct {
		eval       func(interface{}) (want interface{}, result bool)
		input      interface{}
		boolResult bool
	}{
		{testInt, 1, true},
		{testInt, -1, false},
		{testString, "pass", true},
		{testString, "fail", false},
		{testNested, "pass", true},
		{testNested, "fail", false},
	}

	var err error
	for _, provider := range providers {
		if newConf(provider.configure(provider.dat)) != nil {
			t.Fatal("FAILED: %s, %s ", t.Name(), err)
		}
		for _, test := range tests {
			if want, result := test.eval(test.input); result != test.boolResult {
				t.Errorf("FAILED: TestFromFile(conf.json) wanted %v got %v", test.input, want)
			}
		}
	}
}

func TestBadConf(t *testing.T) {

	providers := []struct {
		uri       string
		configure func(string) func(config interface{}) error
	}{
		{"notvalid.json", FromFile},
		{"http://localhost:8080/notvalid.json", FromHTTP},
	}

	for _, provider := range providers {
		err := newConf(provider.configure(provider.uri))
		fmt.Println(err == nil)
		if err == nil {
			t.Errorf("FAILED: TestBadFile(notValid.json) failed to detect bad file %v", conf)
		}
	}
}

func newConf(fn func(config interface{}) error) error {
	return fn(&conf)
}

func testInt(arg interface{}) (want interface{}, result bool) {
	return conf.IntTest, conf.IntTest == arg.(int)
}

func testString(arg interface{}) (want interface{}, result bool) {
	return conf.StringTest, conf.StringTest == arg.(string)
}

func testNested(arg interface{}) (want interface{}, result bool) {
	return conf.Nested.TestNest, conf.StringTest == arg.(string)
}

func readBytes() []byte {
	dat, err := ioutil.ReadFile("./test-config.json")
	if err != nil {
		log.Fatal("FATAL: failed to open test configuration file ", err)
	}
	return dat
}

func startTestServer() chan bool {
	ready := make(chan bool)
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./")))
		go http.ListenAndServe(":8080", nil)
		ready <- true
	}()
	return ready
}

func init() {
	log.Println(<-startTestServer())
}
