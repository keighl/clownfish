
package main

import (
  "testing"
)

func Test_TablePresent_Create(t *testing.T) {
  table := Table{Name: "Test_TablePresent_Create"}

  err := client.TablePresent(table)
  expect(t, err, nil)

  tables, _ := client.TableList()
  expect(t, stringInSlice(table.Name, tables), true)
}

func Test_TablePresent_Create_Fail(t *testing.T) {
  table := Table{
    Name: "Test_TablePresent_Create_Fail####&",
  }

  err := client.TablePresent(table)
  refute(t, err, nil)

  tables, _ := client.TableList()
  expect(t, stringInSlice(table.Name, tables), false)
}

func Test_TablePresent_Exists(t *testing.T) {
  table := Table{
    Name: "Test_TablePresent_Exists",
  }

  err := client.TablePresent(table)
  expect(t, err, nil)

  tables, _ := client.TableList()
  expect(t, stringInSlice(table.Name, tables), true)

  // Look again, shuold not try to create the table
  err = client.TablePresent(table)
  expect(t, err, nil)
}

func Test_TableList_Fail(t *testing.T) {
  c, _ := NewClient(Connection{host, "non-existent-db"})
  _, err := c.TableList()
  refute(t, err, nil)
}

func Test_TableCreateOpts(t *testing.T) {

  table := Table{
    Name: "Test_TableCreateOpts",
  }

  opts := table.TableCreateOpts()
  expect(t, opts.PrimaryKey, nil)
  expect(t, opts.Durability, nil)
  expect(t, opts.DataCenter, nil)

  table = Table{
    Name: "Test_TableCreateOpts",
    PrimaryKey: "cheese",
    Durability: "hard",
    Datacenter: "cheesy_town",
  }

  opts = table.TableCreateOpts()
  expect(t, opts.PrimaryKey, "cheese")
  expect(t, opts.Durability, "hard")
  expect(t, opts.DataCenter, "cheesy_town")

  table = Table{
    Name: "Test_TableCreateOpts",
    Durability: "soft",
  }

  opts = table.TableCreateOpts()
  expect(t, opts.PrimaryKey, nil)
  expect(t, opts.Durability, "soft")
  expect(t, opts.DataCenter, nil)

  table = Table{
    Name: "Test_TableCreateOpts",
    Durability: "cheese",
  }

  opts = table.TableCreateOpts()
  expect(t, opts.PrimaryKey, nil)
  expect(t, opts.Durability, nil)
  expect(t, opts.DataCenter, nil)
}

