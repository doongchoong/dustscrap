package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	From_url       string `json:"from_url"`
	To_url         string `json:"to_url"`
	Mode_url       string `json:"mode_url"`
	Projection_url string `json:"projection_url"`
	Wait_seconds   int64  `json:"wait_seconds"`
	Frames_path    string `json:"frames_path"`
	Time_stamp     bool   `json:"time_stamp"`
}

func (c *Config) Load(confpath string) error {
	dir, _ := filepath.Split(confpath)
	// default value
	c.Wait_seconds = 1
	c.Frames_path = filepath.Join(dir, "frames")
	c.Time_stamp = true

	//open config file
	f, err := os.Open(confpath)
	if err != nil {
		return err
	}
	defer f.Close()

	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	//json
	err = json.Unmarshal(byt, c)
	if err != nil {
		return err
	}

	return nil
}
