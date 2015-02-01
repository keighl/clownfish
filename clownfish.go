package main

import (
  r "github.com/dancannon/gorethink"
  "fmt"
  "errors"
)

// Client manages a rethinkdb connection (scoped to a particular database)
type Client struct {
  session *r.Session
  db string
  LogOutput bool
}

// Log conditionally prints to the standard-out if client.LogOutput is true
func (c *Client) Log(f string) {
  if c.LogOutput {
    printer(f)
  }
}

// NewClient creates a new Client from a Connection
func NewClient(conn Connection) (*Client, error) {
  session, err := r.Connect(r.ConnectOpts{
    Address: conn.Host,
    Database: conn.DB,
  })
  if (err != nil) { return nil, errors.New("Couldn't connect to rethinkdb at "+conn.Host)}
  return &Client{session, conn.DB, false}, nil
}

// Data is a high-level wrapper for YML parsing
type Data struct {
  Conn Connection `yaml:"conn"`
  Tables map[string]Table `yaml:"tables"`
}

// Connection holds info for connecting to a rethinkdb cluster
type Connection struct {
  Host string `yaml:"host"`
  DB string `yaml:"db"`
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

var printer = func (s string) {
  fmt.Println(s)
}
