
package main

import (
  "testing"
  "reflect"
  r "github.com/dancannon/gorethink"
)

var (
  session *r.Session
  host = "localhost:28015"
  db = "clownfish_test"
  client *Client
)

func rebuildDB() {
  r.DbDrop(db).Exec(session)
  r.DbCreate(db).Exec(session)
}

func init() {
  printer = func (s string) {}

  s, err := r.Connect(r.ConnectOpts{
    Address: host,
    Database: db,
  })
  if (err != nil) { panic(err) }
  session = s
  rebuildDB()

  c, _ := NewClient(Connection{host, db})
  client = c
}

func expect(t *testing.T, a interface{}, b interface{}) {
  if a != b {
    t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func refute(t *testing.T, a interface{}, b interface{}) {
  if a == b {
    t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func Test_NewClient(t *testing.T) {
  c, err := NewClient(Connection{host, db})
  expect(t, err, nil)
  expect(t, c.db, "clownfish_test")
}

func Test_NewClient_Error(t *testing.T) {
  _, err := NewClient(Connection{"cheese:8080", db})
  refute(t, err, nil)
}

func Test_ClientLog(t *testing.T) {
  printedString := ""
  printer = func (s string) {
    printedString = s
  }
  c := Client{LogOutput: true}
  c.Log("Client logging...")
  expect(t, printedString, "Client logging...")
}

func Test_ClientLog_Stifle(t *testing.T) {
  receivedPrint := false
  printer = func (s string) {
    receivedPrint = true
  }
  c := Client{LogOutput: false}
  c.Log("Client logging...")
  expect(t, receivedPrint, false)
}

