package main

import (
  r "github.com/dancannon/gorethink"
  "fmt"
)

// Table represents a rethinkdb table definition
type Table struct {
  Name string `yaml:"-"`
  Indicies map[string]Index `yaml:"indices"`
  PrimaryKey string `yaml:"primary_key"`
  Durability string `yaml:"durability"`
  Datacenter string `yaml:"datacenter"`
}

// TableCreateOpts is a factory for gorethink TableCreateOpts
// Ensure durability is either hard or soft
func (x *Table) TableCreateOpts() r.TableCreateOpts {
  opts := r.TableCreateOpts{}

  if (x.PrimaryKey != "") {
    opts.PrimaryKey = x.PrimaryKey
  }

  if (x.Durability == "hard") {
    opts.Durability = "hard"
  }

  if (x.Durability == "soft") {
    opts.Durability = "soft"
  }

  if (x.Datacenter != "") {
    opts.DataCenter = x.Datacenter
  }

  return opts
}

// TablePresent creates a table on client databae if it doesn't exist already
func (c *Client) TablePresent(table Table) error {
  tables, err := c.TableList()
  if (err != nil) { return err }

  if (!stringInSlice(table.Name, tables)) {
    _, err := r.Db(c.db).TableCreate(table.Name, table.TableCreateOpts()).RunWrite(c.session)
    if (err != nil) {
      c.Log(fmt.Sprintf("  %v: table failed to create", table.Name))
      return err
    }
    c.Log(fmt.Sprintf("  %v: table created", table.Name))
  } else {
    c.Log(fmt.Sprintf("  %v:", table.Name))
  }

  return nil
}

// TableList returns a slice of table names on the Client database
func (c *Client) TableList() ([]string, error) {
  res, err := r.Db(c.db).TableList().Run(c.session)
  if (err != nil) { return nil, err }

  tables := []string{}
  res.All(&tables)

  return tables, nil
}