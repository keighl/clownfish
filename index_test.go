
package main

import (
  "testing"
  r "github.com/dancannon/gorethink"
)

func Test_IndexPresent_Simple_Create(t *testing.T) {
  table := Table{Name: "Test_IndexPresent"}
  err := client.TablePresent(table)
  expect(t, err, nil)

  index := Index{
    Name: "Test_IndexPresent_Simple_Create",
  }
  err = client.IndexPresent(table, index)
  expect(t, err, nil)

  indices, _ := client.TableIndexList(table)
  expect(t, stringInSlice(index.Name, indices), true)

  // Do an insert and find
  obj := map[string]string{"Test_IndexPresent_Simple_Create": "cheeseypoofs"}
  _, err = r.Table(table.Name).Insert(obj).RunWrite(session)
  expect(t, err, nil)

  find, err := r.Table(table.Name).
    GetAllByIndex(index.Name, "cheeseypoofs").
    Count().Run(session)
  expect(t, err, nil)
  var count int
  err = find.One(&count)
  expect(t, err, nil)
  expect(t, count, 1)
}

func Test_IndexPresent_Simple_CreateFail(t *testing.T) {
  table := Table{Name: "Test_IndexPresent"}
  err := client.TablePresent(table)
  expect(t, err, nil)

  index := Index{
    Name: "id", // already primary key
  }
  err = client.IndexPresent(table, index)
  refute(t, err, nil)

  indices, _ := client.TableIndexList(table)
  expect(t, stringInSlice(index.Name, indices), false)
}

func Test_IndexPresent_Compound_Create(t *testing.T) {
  table := Table{Name: "Test_IndexPresent"}
  err := client.TablePresent(table)
  expect(t, err, nil)

  index := Index{
    Name: "Test_IndexPresent_Compound_Create",
    Fields: []string{"user_id", "flavor"},
  }
  err = client.IndexPresent(table, index)
  expect(t, err, nil)

  indices, _ := client.TableIndexList(table)
  expect(t, stringInSlice(index.Name, indices), true)


  // Do an insert and find
  obj := map[string]string{"user_id": "12345", "flavor": "vanilla"}
  _, err = r.Table(table.Name).Insert(obj).RunWrite(session)
  expect(t, err, nil)

  find, err := r.Table(table.Name).
    GetAllByIndex(index.Name, []string{"12345", "vanilla"}).
    Count().Run(session)
  expect(t, err, nil)
  var count int
  err = find.One(&count)
  expect(t, err, nil)
  expect(t, count, 1)
}

func Test_IndexPresent_Compound_CreateFail(t *testing.T) {
  table := Table{Name: "Test_IndexPresent"}
  err := client.TablePresent(table)
  expect(t, err, nil)

  index := Index{
    Name: "id", // already a primary key
    Fields: []string{"user_id", "flavor"},
  }
  err = client.IndexPresent(table, index)
  refute(t, err, nil)

  indices, _ := client.TableIndexList(table)
  expect(t, stringInSlice(index.Name, indices), false)
}

func Test_TableIndexList_Fail(t *testing.T) {
  table := Table{Name: "NONEXISTENT_TABLE"}
  _, err := client.TableIndexList(table)
  refute(t, err, nil)
}

func Test_IndexPresent_Exists(t *testing.T) {
  table := Table{Name: "Test_IndexPresent"}
  err := client.TablePresent(table)
  expect(t, err, nil)

  index := Index{
    Name: "Test_IndexPresent_Exists",
  }
  err = client.IndexPresent(table, index)
  expect(t, err, nil)

  // Run it again!
  err = client.IndexPresent(table, index)
  expect(t, err, nil)
}