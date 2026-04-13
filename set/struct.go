/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrDataMustBeStructPointer = errors.New("dataSet must be a struct pointer")

type DataSetStruct struct {
	dataPtr    any           // pointer on struct
	dataStruct reflect.Value // struct
	dataType   reflect.Type
}

func NewDataSetStruct(data any) (*DataSetStruct, error) {
	var dataPtr any
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.Pointer {
		dataType = dataType.Elem()
		dataPtr = data
	} else {
		dataPtr = &data
	}

	if dataType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w, got %T", ErrDataMustBeStructPointer, data)
	}

	return &DataSetStruct{
		dataPtr:    dataPtr,
		dataType:   dataType,
		dataStruct: reflect.Indirect(reflect.ValueOf(data)),
	}, nil
}

//
// func (ds *DataSetStruct) Map() map[string]any {
//	l := ds.dataType.NumField()
//	data := make(map[string]any, l)
//	for i := 0; i < l; i++ {
//		f := ds.dataType.Field(i)
//		name := ds.FieldAliasName(f.Name)
//		data[name] = ds.dataStruct.Field(i).Interface()
//	}
//
//	return data
// }

func (ds *DataSetStruct) FieldValue(name string) (any, error) {
	fieldValue := ds.dataStruct.FieldByName(name)
	if !fieldValue.IsValid() {
		return nil, NewUndefinedFieldError(ds.dataStruct.Interface(), name)
	}

	if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
		return nil, nil
	}

	return fieldValue.Interface(), nil
}

func (ds *DataSetStruct) FieldAliasName(name string) string {
	if field, ok := ds.dataType.FieldByName(name); ok {
		if v, ok := field.Tag.Lookup("json"); ok {
			if name, _, found := strings.Cut(v, ","); found {
				v = name
			}
			name = v
		}
	}

	return name
}

func (ds *DataSetStruct) Name() Name {
	return NameStruct
}

func (ds *DataSetStruct) Data() any {
	return ds.dataPtr
}
