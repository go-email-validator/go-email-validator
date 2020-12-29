package main

import (
	"bytes"
	"encoding/json"
	"github.com/emirpasic/gods/sets/hashset"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	RoleURL      = "https://raw.githubusercontent.com/mixmaxhq/role-based-email-addresses/master/index.js"
	RBEARolePath = "pkg/ev/role/rbea_roles.go"
)

var excludes = hashset.New(
	"asd",
	"asdasd",
	"asdf",
)

func rbeaRolesUpdate(url, path string) {
	rolesResp, err := http.Get(url)
	errPanic(err)
	defer rolesResp.Body.Close()

	var roles = make([]string, 0)
	rolesBytes, err := ioutil.ReadAll(rolesResp.Body)
	errPanic(err)

	rolesBytes = bytes.ReplaceAll(rolesBytes[17:len(rolesBytes)-2], []byte{'\''}, []byte{'"'})
	err = json.Unmarshal(rolesBytes, &roles)
	errPanic(err)

	f, err := os.Create(path)
	errPanic(err)
	defer f.Close()

	f.WriteString(generateRoleCode(roles))
}

func generateRoleCode(roles []string) string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(
		`package role

import "github.com/emirpasic/gods/sets/hashset"

func RBEARoles() []string {
	return rbeaRoles
}

func NewRBEASetRole() SetRole {
	RBEARoles := RBEARoles()
	roles := make([]interface{}, len(RBEARoles))
	for i, role := range RBEARoles {
		roles[i] = role
	}

	return SetRole{hashset.New(roles...)}
}
`)

	strBuilder.WriteString(
		`
var rbeaRoles = []string{
`)
	for _, role := range roles {
		if excludes.Contains(role) {
			continue
		}

		strBuilder.WriteString("\t\"")
		strBuilder.WriteString(role)
		strBuilder.WriteString("\",\n")
	}
	strBuilder.WriteString("}\n")

	return strBuilder.String()
}
