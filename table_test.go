
package main

import (
  "testing"
  r "github.com/dancannon/gorethink"
)

func Test_TablePresent_Create(t *testing.T) {
  table := Table{Name: "Test_TablePresent_Create"}

  err := client.TablePresent(table)
  expect(t, err, nil)

  tables, _ := client.TableList()
  expect(t, stringInSlice(table.Name, tables), true)
}

func Test_TablePresent_TableList_Fail(t *testing.T) {
  c, _ := NewClient(Connection{host, "non-existent-db"})
  table := Table{
    Name: "Test_TablePresent_Create_Fail####&",
  }

  err := c.TablePresent(table)
  refute(t, err, nil)
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

  _, err := r.DB(db).TableCreate(table.Name).RunWrite(client.session)
  expect(t, err, nil)
  client.clearTableList()

  // Run it again!
  err = client.TablePresent(table)
  expect(t, err, nil)
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

func Test_TableAbsent_Success(t *testing.T) {
  table := Table{Name: "Test_TableAbsent_Success"}

  _, err := r.DB(db).TableCreate(table.Name).RunWrite(client.session)
  expect(t, err, nil)
  client.clearTableList()

  err = client.TableAbsent(table)
  expect(t, err, nil)

  tables, _ := client.TableList()
  expect(t, stringInSlice(table.Name, tables), false)
}

func Test_TableAbsent_TableListFail(t *testing.T) {
  table := Table{Name: "Test_TableAbsent_TableListFail"}
  badClient, _ := NewClient(Connection{host, "cheese"})
  err := badClient.TableAbsent(table)
  refute(t, err, nil)
}

func Test_TableAbsent_Fail(t *testing.T) {
  table := Table{Name: "Test_TableAbsent_Fail"}
  // Spoof the table list cache into thinking the table exists already
  client.tableListCache = []string{table.Name}
  err := client.TableAbsent(table)
  refute(t, err, nil)
}

func Test_TableAbsent_Absent(t *testing.T) {
  table := Table{Name: "Test_TableAbsent_Absent"}
  _, err := r.DB(db).TableCreate(table.Name).RunWrite(client.session)
  expect(t, err, nil)

  err = client.TableAbsent(table)
  expect(t, err, nil)

  err = client.TableAbsent(table)
  expect(t, err, nil)
}

