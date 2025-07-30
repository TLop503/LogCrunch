package modules

import (
	"fmt"
	"reflect"
	"regexp"
)

func MetaParse(log string, regex *regexp.Regexp, outputStruct interface{}) error {
	match := regex.FindStringSubmatch(log)
	if match == nil {
		return fmt.Errorf("no match found")
	}

	names := regex.SubexpNames()
	if len(names) != len(match) {
		return fmt.Errorf("mismatch in number of capture groups")
	}

	val := reflect.ValueOf(outputStruct)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("outputStruct must be a non-nil pointer")
	}
	val = val.Elem()
	typ := val.Type()

	for i, name := range names {
		if i == 0 || name == "" {
			continue // skip the full match or unnamed groups
		}

		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			tag := field.Tag.Get("logfield")
			if tag == name {
				if val.Field(j).CanSet() {
					val.Field(j).SetString(match[i])
				}
			}
		}
	}

	return nil
}
