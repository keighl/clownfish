package main

import (
	r "gopkg.in/gorethink/gorethink.v2"
	"reflect"
	"testing"
)

var (
	session *r.Session
	host    = "localhost:28015"
	db      = "clownfish_test"
	client  *Client
)

func rebuildDB() {
	r.DBDrop(db).Exec(session)
	r.DBCreate(db).Exec(session)
}

func init() {
	s, err := r.Connect(r.ConnectOpts{
		Address:  host,
		Database: db,
	})
	if err != nil {
		panic(err)
	}
	session = s
	rebuildDB()

	client, _ = NewClient(Connection{host, db})
	client.LogOutput = false
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

func Test_ClientLog(t *testing.T) {
	c := Client{LogOutput: true}
	c.Log("Client logging...")
}
