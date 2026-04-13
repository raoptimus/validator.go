/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

import "errors"

type DataSetAny struct {
	data any
}

func NewDataSetAny(data any) *DataSetAny {
	return &DataSetAny{
		data: data,
	}
}

func (ds *DataSetAny) FieldValue(_ string) (any, error) {
	return nil, errors.New("not supported")
}

func (ds *DataSetAny) FieldAliasName(name string) string {
	return name
}

func (ds *DataSetAny) Name() Name {
	return NameAny
}

func (ds *DataSetAny) Data() any {
	return ds.data
}
