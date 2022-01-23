package engine_test

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"

	orm "github.com/bhojpur/orm/pkg/engine"
)

func TestScannableSlices(t *testing.T) {
	if err := DB.AutoMigrate(&RecordWithSlice{}).Error; err != nil {
		t.Errorf("Should create table with slice values correctly: %s", err)
	}

	r1 := RecordWithSlice{
		Strings: ExampleStringSlice{"a", "b", "c"},
		Structs: ExampleStructSlice{
			{"name1", "value1"},
			{"name2", "value2"},
		},
	}

	if err := DB.Save(&r1).Error; err != nil {
		t.Errorf("Should save record with slice values")
	}

	var r2 RecordWithSlice

	if err := DB.Find(&r2).Error; err != nil {
		t.Errorf("Should fetch record with slice values")
	}

	if len(r2.Strings) != 3 || r2.Strings[0] != "a" || r2.Strings[1] != "b" || r2.Strings[2] != "c" {
		t.Errorf("Should have serialised and deserialised a string array")
	}

	if len(r2.Structs) != 2 || r2.Structs[0].Name != "name1" || r2.Structs[0].Value != "value1" || r2.Structs[1].Name != "name2" || r2.Structs[1].Value != "value2" {
		t.Errorf("Should have serialised and deserialised a struct array")
	}
}

type RecordWithSlice struct {
	ID      uint64
	Strings ExampleStringSlice `sql:"type:text"`
	Structs ExampleStructSlice `sql:"type:text"`
}

type ExampleStringSlice []string

func (l ExampleStringSlice) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *ExampleStringSlice) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type ExampleStruct struct {
	Name  string
	Value string
}

type ExampleStructSlice []ExampleStruct

func (l ExampleStructSlice) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *ExampleStructSlice) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type ScannerDataType struct {
	Street string `sql:"TYPE:varchar(24)"`
}

func (ScannerDataType) Value() (driver.Value, error) {
	return nil, nil
}

func (*ScannerDataType) Scan(input interface{}) error {
	return nil
}

type ScannerDataTypeTestStruct struct {
	Field1          int
	ScannerDataType *ScannerDataType `sql:"TYPE:json"`
}

type ScannerDataType2 struct {
	Street string `sql:"TYPE:varchar(24)"`
}

func (ScannerDataType2) Value() (driver.Value, error) {
	return nil, nil
}

func (*ScannerDataType2) Scan(input interface{}) error {
	return nil
}

type ScannerDataTypeTestStruct2 struct {
	Field1          int
	ScannerDataType *ScannerDataType2
}

func TestScannerDataType(t *testing.T) {
	scope := orm.Scope{Value: &ScannerDataTypeTestStruct{}}
	if field, ok := scope.FieldByName("ScannerDataType"); ok {
		if DB.Dialect().DataTypeOf(field.StructField) != "json" {
			t.Errorf("data type for scanner is wrong")
		}
	}

	scope = orm.Scope{Value: &ScannerDataTypeTestStruct2{}}
	if field, ok := scope.FieldByName("ScannerDataType"); ok {
		if DB.Dialect().DataTypeOf(field.StructField) != "varchar(24)" {
			t.Errorf("data type for scanner is wrong")
		}
	}
}
