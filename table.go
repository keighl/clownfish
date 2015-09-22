package main

import (
  r "github.com/dancannon/gorethink"
  "fmt"
)

// Table represents a rethinkdb table definition
type Table struct {
  Name string `yaml:"-"`
  Indicies map[string]Index `yaml:"indices"`
  AbsentIndicies []string `yaml:"absent_indices"`
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
    _, err := r.DB(c.db).TableCreate(table.Name, table.TableCreateOpts()).RunWrite(c.session)
    if (err != nil) {
      c.Log(fmt.Sprintf("  + %v ... create failed", table.Name))
      return err
    }
    c.clearTableList()
  }
  c.Log(fmt.Sprintf("  + %v", table.Name))

  return nil
}

// TableAbsent removes a table if it currently exists on the DB
func (c *Client )TableAbsent(table Table) error {
  tables, err := c.TableList()
  if (err != nil) { return err }

  if (stringInSlice(table.Name, tables)) {
    _, err := r.DB(c.db).TableDrop(table.Name).RunWrite(c.session)
    if (err != nil) {
      c.Log(fmt.Sprintf("  - %v ... drop failed", table.Name))
      return err
    }
    c.clearTableList()
  }
  c.Log(fmt.Sprintf("  - %v", table.Name))
  return nil
}

// TableList returns a slice of table names on the Client database
func (c *Client) TableList() ([]string, error) {
  if (len(c.tableListCache) > 0) {
    return c.tableListCache, nil
  }

  res, err := r.DB(c.db).TableList().Run(c.session)
  if (err != nil) { return nil, err }

  c.tableListCache = []string{}
  res.All(&c.tableListCache)

  return c.tableListCache, nil
}

func (c *Client) clearTableList() {
  c.tableListCache = []string{}
}