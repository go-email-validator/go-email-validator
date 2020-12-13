package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	RoleUrl      = "https://raw.githubusercontent.com/mixmaxhq/role-based-email-addresses/master/index.js"
	RBEARolePath = "pkg/ev/role/rbea_roles.go"
)

func main() {
	argsWithoutProg := os.Args[1:]

	var roleUrl = RoleUrl
	if len(argsWithoutProg) > 0 {
		roleUrl = argsWithoutProg[0]
	}

	rolesResp, err := http.Get(roleUrl)
	errPanic(err)
	defer rolesResp.Body.Close()

	var roles = make([]string, 0)
	rolesBytes, err := ioutil.ReadAll(rolesResp.Body)
	errPanic(err)

	rolesBytes = bytes.ReplaceAll(rolesBytes[17:len(rolesBytes)-2], []byte{'\''}, []byte{'"'})
	err = json.Unmarshal(rolesBytes, &roles)
	errPanic(err)

	f, err := os.Create(RBEARolePath)
	errPanic(err)

	defer f.Close()

	f.WriteString(generateCode(roles))

	println(roles)
	println(err)
}

func generateCode(roles []string) string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(
		`package role

import "github.com/emirpasic/gods/sets/hashset"

var RBEARoles = []string{
`)
	for _, role := range roles {
		strBuilder.WriteString("\t\"")
		strBuilder.WriteString(role)
		strBuilder.WriteString("\",\n")
	}
	strBuilder.WriteString("}\n")

	strBuilder.WriteString(`
func NewRBEASetRole() SetRole {
	roles := make([]interface{}, len(RBEARoles))
	for i, role := range RBEARoles {
		roles[i] = role
	}

	return SetRole{hashset.New(roles...)}
}
`)

	return strBuilder.String()
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
