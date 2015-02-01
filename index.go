package main

import (
  r "github.com/dancannon/gorethink"
  "fmt"
)

// Index represents a rethinkdb index definition
type Index struct {
  // The name of the index
  Name string `yaml:"-"`
  // Multiple fields for a compound index
  Fields []string `yaml:"fields"`
  Multi bool `yaml:"multi"`
  Geo bool `yaml:"geo"`
}

// IndexCreateOpts is a factory for gorethink IndexCreateOpts
func (x *Index) IndexCreateOpts() r.IndexCreateOpts {
  return r.IndexCreateOpts{Multi: x.Multi, Geo: x.Geo}
}

// IndexPresent xreates an index on the table if it doesn't exist already
func (c *Client) IndexPresent(table Table, index Index) error {
  indices, err := c.TableIndexList(table)
  if (err != nil) { return err }

  if (!stringInSlice(index.Name, indices)) {
    if (len(index.Fields) == 0) {
      _, err = r.Table(table.Name).IndexCreate(index.Name, index.IndexCreateOpts()).RunWrite(c.session)
    } else {
      _, err = r.Table(table.Name).IndexCreateFunc(index.Name, func(row r.Term) interface{} {
        ar := []interface{}{}
        for _, field := range index.Fields {
          ar = append(ar, row.Field(field))
        }
        return ar
      }, index.IndexCreateOpts()).RunWrite(c.session)
    }

    if (err != nil) {
      c.Log(fmt.Sprintf("      * %s: failed to create", index.Name))
      return err
    }
    c.Log(fmt.Sprintf("      * %s: index created", index.Name))
  } else {
    c.Log(fmt.Sprintf("      * %s", index.Name))
  }

  return nil
}

// TableIndexList returns a slice of index names on the table
func (c *Client) TableIndexList(table Table) ([]string, error) {
  res, err := r.Table(table.Name).IndexList().Run(c.session)
  if (err != nil) { return nil, err }

  indices := []string{}
  res.All(&indices)
  return indices, nil
}
