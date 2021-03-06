package main

import (
	"errors"
	"fmt"
	r "gopkg.in/gorethink/gorethink.v2"
)

// Client manages a rethinkdb connection (scoped to a particular database)
type Client struct {
	session        *r.Session
	db             string
	LogOutput      bool
	indexListCache map[string][]string
	tableListCache []string
}

// Log conditionally prints to the standard-out if client.LogOutput is true
func (c *Client) Log(f string) {
	if c.LogOutput {
		fmt.Println(f)
	}
}

// NewClient creates a new Client from a Connection
func NewClient(conn Connection) (*Client, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  conn.Host,
		Database: conn.DB,
	})
	if err != nil {
		return nil, errors.New("Couldn't connect to rethinkdb at " + conn.Host)
	}
	return &Client{
		session:        session,
		db:             conn.DB,
		LogOutput:      false,
		indexListCache: map[string][]string{},
		tableListCache: []string{},
	}, nil
}

// Data is a high-level wrapper for YML parsing
type Data struct {
	Conn         Connection       `yaml:"conn"`
	Tables       map[string]Table `yaml:"tables"`
	AbsentTables []string         `yaml:"absent_tables"`
}

// Connection holds info for connecting to a rethinkdb cluster
type Connection struct {
	Host string `yaml:"host"`
	DB   string `yaml:"db"`
}

// UTILS ////////////////////

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
