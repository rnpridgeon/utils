package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config map[string]interface{}

// Returns function for reading in config file then passing bytes to be marshaled into application provided config object
func FromFile(path string) func(config interface{}) (err error) {
	return func(config interface{}) (err error) {
		var dat []byte
		if dat, err = ioutil.ReadFile(path); err == nil {
			return json.Unmarshal(dat, &config)
		}
		return err
	}
}

// Returns function which unmarshal bytes into application provided config object
func FromBytes(dat []byte) func(config interface{}) (err error) {
	return func(config interface{}) (err error) {
		return json.Unmarshal(dat, &config)
	}
}

// TODO: add test
// TODO: will require reflection to keep things working the way I want
// Returns function for providing existing config object to application
//func FromConfig(conf interface{}) func(config interface{}) (err error) {
//	return func(config interface{}) (err error) {
//		return nil
//	}
//}

// TODO: add Vault variant
// Returns function for fetching and processing bytes to be marshaled into application provided config object
func FromHTTP(uri string) func(config interface{}) (err error) {
	return func(config interface{}) (err error) {
		var resp *http.Response
		if resp, err = http.Get(uri); err == nil && resp.StatusCode == http.StatusOK {
			return json.NewDecoder(resp.Body).Decode(&config)
		}
		if err != nil {
			return err
		} else {
			return errors.New(fmt.Sprintf("Non 200 status code returned: %s", resp.Status))
		}
	}
}
