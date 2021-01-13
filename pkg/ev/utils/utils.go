package utils

import (
	"github.com/joho/godotenv"
	"log"
	"reflect"
)

// RangeLen returns length of interface{}
func RangeLen(i interface{}) int {
	if i == nil {
		return 0
	}

	return reflect.ValueOf(i).Len()
}

// Errs forms []error from errors...
func Errs(errs ...error) []error {
	if errs == nil || len(errs) == 1 && errs[0] == nil {
		return nil
	}

	return errs
}

// LoadEnv loads
func LoadEnv(env string) {
	filenames := make([]string, 0)

	if env != "" {
		filenames = append(filenames, env)
	}

	if err := godotenv.Load(filenames...); err != nil {
		log.Print("No .env file found")
	}
}

// StructName returns name of structure
func StructName(strct interface{}) string {
	return reflect.ValueOf(strct).Type().String()
}
