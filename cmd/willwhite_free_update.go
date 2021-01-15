package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

const (
	freeURL           = "https://raw.githubusercontent.com/willwhite/freemail/master/data/free.txt"
	willWhiteFreePath = "pkg/ev/free/willwhite_free.go"
)

func willwhiteFreeUpdate(url, path string) {
	freeResp, err := http.Get(url)
	errPanic(err)
	defer freeResp.Body.Close()

	var freeEmails = make([]string, 0)
	bufReader := bufio.NewReader(freeResp.Body)

	for {
		freeEmail, _, err := bufReader.ReadLine()
		if len(freeEmail) > 0 {
			freeEmails = append(freeEmails, string(freeEmail))
		}
		if err != nil {
			break
		}
	}

	f, err := os.Create(path)
	errPanic(err)
	defer f.Close()

	f.WriteString(generateFreeCode(freeEmails))
}

func generateFreeCode(freeEmails []string) string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(
		`package free

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"strings"
)

// WillWhiteFree returns the list of free domains
func WillWhiteFree() []string {
	return willWhiteFree
}

// NewWillWhiteSetFree forms contains.InSet from list of free domains (https://github.com/willwhite/freemail/blob/master/data/free.txt)
func NewWillWhiteSetFree() contains.InSet {
	WillWhiteFree := WillWhiteFree()
	freeEmails := make([]interface{}, len(WillWhiteFree))
	for i, freeEmail := range WillWhiteFree {
		freeEmails[i] = strings.ToLower(freeEmail)
	}

	return contains.NewSet(hashset.New(freeEmails...))
}
`)
	strBuilder.WriteString(`
var willWhiteFree = []string{
`)
	for _, freeEmail := range freeEmails {
		strBuilder.WriteString("\t\"")
		strBuilder.WriteString(freeEmail)
		strBuilder.WriteString("\",\n")
	}
	strBuilder.WriteString("}\n")

	return strBuilder.String()
}
