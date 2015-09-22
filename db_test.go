
package main

import (
  "testing"
  r "github.com/dancannon/gorethink"
)

func Test_DBPresent_Create(t *testing.T) {
  client, _  := NewClient(Connection{host, "Test_DBPresent_Create"})

  dbs := client.DBList()
  expect(t, stringInSlice(client.db, dbs), false)

  err := client.DBPresent()
  expect(t, err, nil)

  dbs = client.DBList()
  expect(t, stringInSlice(client.db, dbs), true)

  _, err = r.DBDrop(client.db).Run(session)
  expect(t, err, nil)
}

func Test_DBPresent_Create_Fail(t *testing.T) {
  client, _  := NewClient(Connection{host, "Test_DBPresent_Create_Fail#####"})

  dbs := client.DBList()
  expect(t, stringInSlice(client.db, dbs), false)

  err := client.DBPresent()
  refute(t, err, nil)

  dbs = client.DBList()
  expect(t, stringInSlice(client.db, dbs), false)
}

func Test_DBPresent_Exists(t *testing.T) {
  client, _  := NewClient(Connection{host, "Test_DBPresent_Exists"})

  dbs := client.DBList()
  expect(t, stringInSlice(client.db, dbs), false)

  err := client.DBPresent()
  expect(t, err, nil)

  dbs = client.DBList()
  expect(t, stringInSlice(client.db, dbs), true)

  err = client.DBPresent()
  expect(t, err, nil)

  _, err = r.DBDrop(client.db).Run(session)
  expect(t, err, nil)
}
