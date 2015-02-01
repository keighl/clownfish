package main

import (
  "testing"
)

func Test_ParseYMLFile_Fail_NoFile(t *testing.T) {
  err := ParseYMLFile("test/nonexistent.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Fail_BadYAML(t *testing.T) {
  err := ParseYMLFile("test/bad_yaml.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Fail_BadConnection(t *testing.T) {
  err := ParseYMLFile("test/bad_connection.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Fail_DBPresentError(t *testing.T) {
  err := ParseYMLFile("test/db_present_error.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Fail_TablePresentError(t *testing.T) {
  err := ParseYMLFile("test/table_present_error.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Fail_IndexPresentError(t *testing.T) {
  err := ParseYMLFile("test/index_present_error.yml")
  refute(t, err, nil)
}

func Test_ParseYMLFile_Success(t *testing.T) {
  err := ParseYMLFile("test/success.yml")
  expect(t, err, nil)
}