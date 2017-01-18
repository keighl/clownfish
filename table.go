package main

import (
	"fmt"
	r "gopkg.in/gorethink/gorethink.v2"
)

// Table represents a rethinkdb table definition
type Table struct {
	Name           string           `yaml:"-"`
	Indicies       map[string]Index `yaml:"indices"`
	AbsentIndicies []string         `yaml:"absent_indices"`

	// TableCreateOpts
	PrimaryKey           interface{} `yaml:primary_key`
	Durability           interface{} `yaml:durability`
	Shards               interface{} `yaml:shards`
	Replicas             interface{} `yaml:replicas`
	PrimaryReplicaTag    interface{} `yaml:primary_replica_tag`
	NonVotingReplicaTags interface{} `yaml:nonvoting_replica_tags`
}

// TableCreateOpts is a factory for gorethink TableCreateOpts
func (x *Table) TableCreateOpts() r.TableCreateOpts {
	opts := r.TableCreateOpts{
		PrimaryKey:           x.PrimaryKey,
		Durability:           x.Durability,
		Shards:               x.Shards,
		Replicas:             x.Replicas,
		PrimaryReplicaTag:    x.PrimaryReplicaTag,
		NonVotingReplicaTags: x.NonVotingReplicaTags,
	}

	return opts
}

// TablePresent creates a table on client databae if it doesn't exist already
func (c *Client) TablePresent(table Table) error {
	tables, err := c.TableList()
	if err != nil {
		return err
	}

	if !stringInSlice(table.Name, tables) {
		_, err := r.DB(c.db).TableCreate(table.Name, table.TableCreateOpts()).RunWrite(c.session)
		if err != nil {
			c.Log(fmt.Sprintf("  + %v ... create failed", table.Name))
			return err
		}
		c.clearTableList()
	}
	c.Log(fmt.Sprintf("  + %v", table.Name))

	return nil
}

// TableAbsent removes a table if it currently exists on the DB
func (c *Client) TableAbsent(table Table) error {
	tables, err := c.TableList()
	if err != nil {
		return err
	}

	if stringInSlice(table.Name, tables) {
		_, err := r.DB(c.db).TableDrop(table.Name).RunWrite(c.session)
		if err != nil {
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
	if len(c.tableListCache) > 0 {
		return c.tableListCache, nil
	}

	res, err := r.DB(c.db).TableList().Run(c.session)
	if err != nil {
		return nil, err
	}

	c.tableListCache = []string{}
	res.All(&c.tableListCache)

	return c.tableListCache, nil
}

func (c *Client) clearTableList() {
	c.tableListCache = []string{}
}
