package main

import (
  "os"
  "github.com/codegangsta/cli"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "fmt"
)

var defaultFileName = "clownfish.yml"

func main() {
  newCLI().Run(os.Args)
}

func newCLI() *cli.App {
  app := cli.NewApp()
  app.Name = "clownfish"
  app.Usage = `YAML based management tool for RethinkDB tables and indices.

   $ clownfish
   $ clownfish db_config.yml`
  app.Version = "0.1.0"
  app.Action = cliAction
  return app
}

func cliAction(c *cli.Context) {
  file := defaultFileName

  if len(c.Args()) > 0 {
    file = c.Args()[0]
  }

  err := ParseYMLFile(file)
  if (err != nil) {
    fmt.Println(err.Error())
    os.Exit(1)
  }
}

// ParseYMLFile opens a YML file, and configures the specified rethink database accordingly
func ParseYMLFile(file string) error {

  data, err := ioutil.ReadFile(file)
  if (err != nil) { return err }

  d := Data{}
  err = yaml.Unmarshal([]byte(data), &d)
  if (err != nil) { return err }

  client, err := NewClient(d.Conn)
  if (err != nil) { return err }
  client.LogOutput = true

  err = client.DBPresent()
  if (err != nil) { return err }

  for name, table := range d.Tables {

    // create table
    table.Name = name
    err := client.TablePresent(table)
    if (err != nil) { return err }

    // create indices
    for indexName, index := range table.Indicies {
      index.Name = indexName
      err := client.IndexPresent(table, index)
      if (err != nil) { return err }
    }

    for _, indexName := range table.AbsentIndicies {
      index := Index{Name: indexName}
      err := client.IndexAbsent(table, index)
      if (err != nil) { return err }
    }
  }

  // remove tables
  for _, tableName := range d.AbsentTables {
    table := Table{Name: tableName}
    err := client.TableAbsent(table)
    if (err != nil) { return err }
  }

  return nil
}
