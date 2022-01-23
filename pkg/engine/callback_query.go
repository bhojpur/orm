package engine

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
	"errors"
	"fmt"
	"reflect"
)

// Define callbacks for querying
func init() {
	DefaultCallback.Query().Register("orm:query", queryCallback)
	DefaultCallback.Query().Register("orm:preload", preloadCallback)
	DefaultCallback.Query().Register("orm:after_query", afterQueryCallback)
}

// queryCallback used to query data from database
func queryCallback(scope *Scope) {
	if _, skip := scope.InstanceGet("orm:skip_query_callback"); skip {
		return
	}

	//we are only preloading relations, dont touch base model
	if _, skip := scope.InstanceGet("orm:only_preload"); skip {
		return
	}

	defer scope.trace(NowFunc())

	var (
		isSlice, isPtr bool
		resultType     reflect.Type
		results        = scope.IndirectValue()
	)

	if orderBy, ok := scope.Get("orm:order_by_primary_key"); ok {
		if primaryField := scope.PrimaryField(); primaryField != nil {
			scope.Search.Order(fmt.Sprintf("%v.%v %v", scope.QuotedTableName(), scope.Quote(primaryField.DBName), orderBy))
		}
	}

	if value, ok := scope.Get("orm:query_destination"); ok {
		results = indirect(reflect.ValueOf(value))
	}

	if kind := results.Kind(); kind == reflect.Slice {
		isSlice = true
		resultType = results.Type().Elem()
		results.Set(reflect.MakeSlice(results.Type(), 0, 0))

		if resultType.Kind() == reflect.Ptr {
			isPtr = true
			resultType = resultType.Elem()
		}
	} else if kind != reflect.Struct {
		scope.Err(errors.New("unsupported destination, should be slice or struct"))
		return
	}

	scope.prepareQuerySQL()

	if !scope.HasError() {
		scope.db.RowsAffected = 0

		if str, ok := scope.Get("orm:query_hint"); ok {
			scope.SQL = fmt.Sprint(str) + scope.SQL
		}

		if rows, err := scope.SQLDB().Query(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
			defer rows.Close()

			columns, _ := rows.Columns()
			for rows.Next() {
				scope.db.RowsAffected++

				elem := results
				if isSlice {
					elem = reflect.New(resultType).Elem()
				}

				scope.scan(rows, columns, scope.New(elem.Addr().Interface()).Fields())

				if isSlice {
					if isPtr {
						results.Set(reflect.Append(results, elem.Addr()))
					} else {
						results.Set(reflect.Append(results, elem))
					}
				}
			}

			if err := rows.Err(); err != nil {
				scope.Err(err)
			} else if scope.db.RowsAffected == 0 && !isSlice {
				scope.Err(ErrRecordNotFound)
			}
		}
	}
}

// afterQueryCallback will invoke `AfterFind` method after querying
func afterQueryCallback(scope *Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterFind")
	}
}
