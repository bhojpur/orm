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
	"fmt"
	"reflect"
	"testing"
)

func TestCloneSearch(t *testing.T) {
	s := new(search)
	s.Where("name = ?", "shashi").Order("name").Attrs("name", "shashi").Select("name, age")

	s1 := s.clone()
	s1.Where("age = ?", 20).Order("age").Attrs("email", "a@e.org").Select("email")

	if reflect.DeepEqual(s.whereConditions, s1.whereConditions) {
		t.Errorf("Where should be copied")
	}

	if reflect.DeepEqual(s.orders, s1.orders) {
		t.Errorf("Order should be copied")
	}

	if reflect.DeepEqual(s.initAttrs, s1.initAttrs) {
		t.Errorf("InitAttrs should be copied")
	}

	if reflect.DeepEqual(s.Select, s1.Select) {
		t.Errorf("selectStr should be copied")
	}
}

func TestWhereCloneCorruption(t *testing.T) {
	for whereCount := 1; whereCount <= 8; whereCount++ {
		t.Run(fmt.Sprintf("w=%d", whereCount), func(t *testing.T) {
			s := new(search)
			for w := 0; w < whereCount; w++ {
				s = s.clone().Where(fmt.Sprintf("w%d = ?", w), fmt.Sprintf("value%d", w))
			}
			if len(s.whereConditions) != whereCount {
				t.Errorf("s: where count should be %d", whereCount)
			}

			q1 := s.clone().Where("finalThing = ?", "THING1")
			q2 := s.clone().Where("finalThing = ?", "THING2")

			if reflect.DeepEqual(q1.whereConditions, q2.whereConditions) {
				t.Errorf("Where conditions should be different")
			}
		})
	}
}
