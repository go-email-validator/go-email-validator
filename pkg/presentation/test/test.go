package test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// DefaultDepFixtureFile is a name of test file
const DefaultDepFixtureFile = "dep_fixture_test.json"

// DepPresentations returns structs from json test file
func DepPresentations(t *testing.T, result interface{}, fp string) {
	if fp == "" {
		fp = DefaultDepFixtureFile
	}

	fp, err := filepath.Abs(fp)
	require.Nil(t, err)
	jsonFile, err := os.Open(fp)
	require.Nil(t, err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.Nil(t, err)

	err = json.Unmarshal(byteValue, &result)
	require.Nil(t, err)
}
