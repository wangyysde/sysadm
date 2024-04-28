package main

import (
	"fmt"
	"reflect"
	"strings"
)

func GetFeildValueByName(data any, fieldName string) (interface{}, error) {
	dT := reflect.TypeOf(data)
	dV := reflect.ValueOf(data)
	if dT.Kind() == reflect.Pointer {
		dT = dT.Elem()
		dV = dV.Elem()
	}
	if dT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("we can not only get a struct field value")
	}

	for i := 0; i < dT.NumField(); i++ {
		field := dT.Field(i)
		fmt.Printf("field name is: %s\n", field.Name)
		if field.Name == fieldName {
			fieldT := field.Type
			v := reflect.ValueOf(fieldT)
			fmt.Printf("create new value is: %s\n", v)
			fieldValue := dV.Field(i)
			return fieldValue.Interface(), nil
		}
	}

	return nil, fmt.Errorf("field named %s was not found in struct %s", fieldName, dT.Name())
}

func setFieldValue(data, value any, fieldName string) error {
	dT := reflect.TypeOf(data)
	if dT.Kind() != reflect.Pointer || dT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer point to a struct")
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		return fmt.Errorf("field name must not empty")
	}
	dTElem := dT.Elem()
	dV := reflect.ValueOf(data).Elem()
	for i := 0; i < dTElem.NumField(); i++ {
		field := dTElem.Field(i)
		if field.Name == fieldName {
			fieldType := field.Type
			valueType := reflect.TypeOf(value)
			if fieldType != valueType {
				return fmt.Errorf("the type of value is not equal to the type of the field of the data")
			}
			tmpValue := reflect.ValueOf(value)
			dV.Field(i).Set(tmpValue)
		}
	}

	return nil
}
