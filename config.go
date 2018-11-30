package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	ucl "github.com/nahanni/go-ucl"
)

// Config holds the runtime configuration which is expected to be
// read from a UCL formatted file
type Config struct {
	Broker string `json:"broker.address"`
	Topic  string `json:"topic"`

	SASL struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"sasl,omitempty"`

	WebServer struct {
		Path          string `json:"path"`
		Port          string `json:"port"`
		BasicAuthFile string `json:"basicauthfile"`

		TLS struct {
			Cert string `json:"certfile"`
			Key  string `json:"keyfile"`
		} `json:"tls,omitempty"`
	} `json:"webserver"`
}

//FromFile sets config based on file (fname) content
func (c *Config) FromFile(fname string) error {
	var (
		file, uclJSON []byte
		err           error
		fileBytes     *bytes.Buffer
		parser        *ucl.Parser
		uclData       map[string]interface{}
	)
	if fname, err = filepath.Abs(fname); err != nil {
		return err
	}
	if fname, err = filepath.EvalSymlinks(fname); err != nil {
		return err
	}
	if file, err = ioutil.ReadFile(fname); err != nil {
		return err
	}

	fileBytes = bytes.NewBuffer(file)
	parser = ucl.NewParser(fileBytes)
	if uclData, err = parser.Ucl(); err != nil {
		return err
	}
	if uclJSON, err = json.Marshal(uclData); err != nil {
		return err
	}
	return json.Unmarshal(uclJSON, &c)
}
