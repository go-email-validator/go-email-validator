package utils

import (
	"github.com/joho/godotenv"
	"log"
	"reflect"
	"time"
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

var (
	defaultString string
	defaultInt    int
)

// DefaultString return defaultVal if val is empty
func DefaultString(val string, defaultVal string) string {
	if val == defaultString {
		return defaultVal
	}

	return val
}

// DefaultInt return defaultVal if val is empty
func DefaultInt(val int, defaultVal int) int {
	if val == defaultInt {
		return defaultVal
	}

	return val
}

// DefaultDuration return defaultVal if val is empty fpr time.Duration
func DefaultDuration(val time.Duration, defaultVal time.Duration) time.Duration {
	if val == 0 {
		return defaultVal
	}

	return val
}

// DefaultInterface return defaultVal if val is empty
func DefaultInterface(val interface{}, defaultVal interface{}) interface{} {
	if val == nil {
		return defaultVal
	}

	return val
}
