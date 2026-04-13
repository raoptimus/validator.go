/**
 * This file is part of the raoptimus/validator.go library
 *
 * @copyright Copyright (c) Evgeniy Urvantsev
 * @license https://github.com/raoptimus/validator.go/blob/master/LICENSE.md
 * @link https://github.com/raoptimus/validator.go
 */
package set

type DataSetMap struct {
	data map[string]any
}

func NewDataSetMap(data map[string]any) *DataSetMap {
	return &DataSetMap{
		data: data,
	}
}

func (ds *DataSetMap) FieldValue(name string) (any, error) {
	v, ok := ds.data[name]
	if !ok {
		return nil, NewUndefinedFieldError(ds.data, name)
	}

	return v, nil
}

func (ds *DataSetMap) FieldAliasName(name string) string {
	return name
}

func (ds *DataSetMap) Name() Name {
	return NameMap
}

func (ds *DataSetMap) Data() any {
	return ds.data
}
