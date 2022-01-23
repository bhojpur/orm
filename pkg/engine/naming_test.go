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
	"testing"

	orm "github.com/bhojpur/orm/pkg/engine"
)

func TestTheNamingStrategy(t *testing.T) {

	cases := []struct {
		name     string
		namer    orm.Namer
		expected string
	}{
		{name: "auth", expected: "auth", namer: orm.TheNamingStrategy.DB},
		{name: "userRestrictions", expected: "user_restrictions", namer: orm.TheNamingStrategy.Table},
		{name: "clientID", expected: "client_id", namer: orm.TheNamingStrategy.Column},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.namer(c.name)
			if result != c.expected {
				t.Errorf("error in naming strategy. expected: %v got :%v\n", c.expected, result)
			}
		})
	}

}

func TestNamingStrategy(t *testing.T) {

	dbNameNS := func(name string) string {
		return "db_" + name
	}
	tableNameNS := func(name string) string {
		return "tbl_" + name
	}
	columnNameNS := func(name string) string {
		return "col_" + name
	}

	ns := &orm.NamingStrategy{
		DB:     dbNameNS,
		Table:  tableNameNS,
		Column: columnNameNS,
	}

	cases := []struct {
		name     string
		namer    orm.Namer
		expected string
	}{
		{name: "auth", expected: "db_auth", namer: ns.DB},
		{name: "user", expected: "tbl_user", namer: ns.Table},
		{name: "password", expected: "col_password", namer: ns.Column},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.namer(c.name)
			if result != c.expected {
				t.Errorf("error in naming strategy. expected: %v got :%v\n", c.expected, result)
			}
		})
	}

}
