package utils

import (
	"github.com/joho/godotenv"
	"log"
	"reflect"
)

func RangeLen(i interface{}) int {
	if i == nil {
		return 0
	}

	return reflect.ValueOf(i).Len()
}

func Errs(errs ...error) []error {
	if errs == nil || len(errs) == 1 && errs[0] == nil {
		return nil
	}

	return errs
}

type Err struct {
	s string
}

func (e Err) Error() string {
	return e.s
}

func LoadEnv(env string) {
	filenames := make([]string, 0)

	if env != "" {
		filenames = append(filenames, env)
	}

	if err := godotenv.Load(filenames...); err != nil {
		log.Print("No .env file found")
	}
}
