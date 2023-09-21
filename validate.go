package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	RequiredTag = "required"
	NotNull     = "not_null"
	NotZero     = "not_zero"
)

type Options struct {
	NoRedundant bool `json:"no_redundant,omitempty,not_null"` // cannot have fields not defined in struct
}

func WithNoRedundant(opt *Options) {
	opt.NoRedundant = true
}

type Option func(*Options)

func Validate(data []byte, obj interface{}, opts ...Option) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("invalid json")
	}

	var opt Options
	for _, fn := range opts {
		fn(&opt)
	}
	return validateObject(m, reflect.TypeOf(obj), &opt)
}

func validateObject(m map[string]interface{}, typ reflect.Type, opt *Options) error {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("invalid object type")
	}

	tag2Field := map[string]*reflect.StructField{}
	tag2Tags := map[string][]string{}

	// if field is required, but not exists in m
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		tags := strings.Split(f.Tag.Get("json"), ",")
		jsonName := tags[0]
		if jsonName == "" {
			jsonName = f.Name
		}
		tag2Field[jsonName] = &f
		tag2Tags[jsonName] = tags[1:]

		for _, tag := range tags[1:] {
			tag = strings.TrimSpace(tag)
			switch tag {
			case RequiredTag:
				if _, ok := m[jsonName]; !ok {
					return fmt.Errorf("field %s is required", f.Name)
				}
			}
		}
	}

	for k, v := range m {
		f, ok := tag2Field[k]
		if !ok {
			if opt.NoRedundant {
				return fmt.Errorf("redundant field %s", k)
			}
			continue
		}

		for _, tag := range tag2Tags[k] {
			tag = strings.TrimSpace(tag)
			switch tag {
			case NotNull:
				if v == nil {
					return fmt.Errorf("field %s cannot be null", k)
				}
			case NotZero:
				if v == nil || reflect.ValueOf(v).IsZero() {
					return fmt.Errorf("field %s cannot be zero", k)
				}
			}
		}

		// if f is struct
		switch f.Type.Kind() {
		case reflect.Struct, reflect.Pointer:
			if v == nil {
				break
			}
			sub, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid field %s, need map, but got %T", k, v)
			}

			if err := validateObject(sub, f.Type, opt); err != nil {
				return err
			}
		}
	}

	return nil
}
