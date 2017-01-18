package main

import (
	r "gopkg.in/gorethink/gorethink.v2"
	"testing"
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
