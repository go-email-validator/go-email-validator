package presentation_test

import (
	"encoding/json"
	openapi "github.com/go-email-validator/go-ev-presenters/pkg/api/v1/go"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const DefaultDepFixtureFile = "dep_fixture_test.json"

func TestDepPresentations(t *testing.T, result interface{}, fp string) {
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

type unmarshalString []byte

func (v *unmarshalString) UnmarshalJSON(data []byte) error {
	*v = data
	return nil
}

func TestEmailResponses(
	t *testing.T,
	filepath,
	path string,
	factory func(data []byte) *openapi.EmailResponse,
) []*openapi.EmailResponse {
	if filepath == "" {
		filepath = DefaultDepFixtureFile
	}

	if path == "" {
		path = "@this"
	}

	jsonFile, err := os.Open(filepath)
	require.Nil(t, err, filepath)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.Nil(t, err, filepath)

	value := gjson.Get(string(byteValue), path)
	jsonResponses := make([]unmarshalString, 0)
	err = json.Unmarshal([]byte(value.String()), &jsonResponses)
	require.Nil(t, err, filepath)

	responses := make([]*openapi.EmailResponse, len(jsonResponses))
	for index, jsonResponse := range jsonResponses {
		responses[index] = factory(jsonResponse)
	}

	return responses
}
