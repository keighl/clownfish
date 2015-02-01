package main

import (
  r "github.com/dancannon/gorethink"
  "fmt"
)

// DBPresent creates a database in the rethink cluster if it doesn't exist already
func (c *Client) DBPresent() error {
  if (!stringInSlice(c.db, c.DBList())) {

    _, err := r.DbCreate(c.db).RunWrite(c.session)
    if (err != nil) {
      c.Log(fmt.Sprintf(":%v DB failed to create", c.db))
      return err
    }
    c.Log(fmt.Sprintf("%v: DB created", c.db))
  } else {
    c.Log(fmt.Sprintf("%v:", c.db))
  }

  return nil
}

// DBList returns a slice of cluster database names
func (c *Client) DBList() []string {
  res, _ := r.DbList().Run(c.session)
  dbs := []string{}
  res.All(&dbs)
  return dbs
}