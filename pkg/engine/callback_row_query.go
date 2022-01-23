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
	"database/sql"
	"fmt"
)

// Define callbacks for row query
func init() {
	DefaultCallback.RowQuery().Register("orm:row_query", rowQueryCallback)
}

type RowQueryResult struct {
	Row *sql.Row
}

type RowsQueryResult struct {
	Rows  *sql.Rows
	Error error
}

// queryCallback used to query data from database
func rowQueryCallback(scope *Scope) {
	if result, ok := scope.InstanceGet("row_query_result"); ok {
		scope.prepareQuerySQL()

		if str, ok := scope.Get("orm:query_hint"); ok {
			scope.SQL = fmt.Sprint(str) + scope.SQL
		}

		if rowResult, ok := result.(*RowQueryResult); ok {
			rowResult.Row = scope.SQLDB().QueryRow(scope.SQL, scope.SQLVars...)
		} else if rowsResult, ok := result.(*RowsQueryResult); ok {
			rowsResult.Rows, rowsResult.Error = scope.SQLDB().Query(scope.SQL, scope.SQLVars...)
		}
	}
}
