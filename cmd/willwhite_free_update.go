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

import "github.com/emirpasic/gods/sets/hashset"

func WillWhiteFree() []string {
	return willWhiteFree
}

func NewWillWhiteSetFree() SetFree {
	WillWhiteFree := WillWhiteFree()
	freeEmails := make([]interface{}, len(WillWhiteFree))
	for i, freeEmail := range WillWhiteFree {
		freeEmails[i] = freeEmail
	}

	return SetFree{hashset.New(freeEmails...)}
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
