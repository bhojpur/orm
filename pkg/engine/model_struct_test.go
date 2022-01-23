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
	"sync"
	"testing"

	orm "github.com/bhojpur/orm/pkg/engine"
)

type ModelA struct {
	orm.Model
	Name string

	ModelCs []ModelC `orm:"foreignkey:OtherAID"`
}

type ModelB struct {
	orm.Model
	Name string

	ModelCs []ModelC `orm:"foreignkey:OtherBID"`
}

type ModelC struct {
	orm.Model
	Name string

	OtherAID uint64
	OtherA   *ModelA `orm:"foreignkey:OtherAID"`
	OtherBID uint64
	OtherB   *ModelB `orm:"foreignkey:OtherBID"`
}

type RequestModel struct {
	Name     string
	Children []ChildModel `orm:"foreignkey:ParentID"`
}

type ChildModel struct {
	ID       string
	ParentID string
	Name     string
}

type ResponseModel struct {
	orm.Model
	RequestModel
}

// This test will try to cause a race condition on the model's foreignkey metadata
func TestModelStructRaceSameModel(t *testing.T) {
	// use a WaitGroup to execute as much in-sync as possible
	// it's more likely to hit a race condition than without
	n := 32
	start := sync.WaitGroup{}
	start.Add(n)

	// use another WaitGroup to know when the test is done
	done := sync.WaitGroup{}
	done.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			start.Wait()

			// call GetStructFields, this had a race condition before we fixed it
			DB.NewScope(&ModelA{}).GetStructFields()

			done.Done()
		}()

		start.Done()
	}

	done.Wait()
}

// This test will try to cause a race condition on the model's foreignkey metadata
func TestModelStructRaceDifferentModel(t *testing.T) {
	// use a WaitGroup to execute as much in-sync as possible
	// it's more likely to hit a race condition than without
	n := 32
	start := sync.WaitGroup{}
	start.Add(n)

	// use another WaitGroup to know when the test is done
	done := sync.WaitGroup{}
	done.Add(n)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			start.Wait()

			// call GetStructFields, this had a race condition before we fixed it
			if i%2 == 0 {
				DB.NewScope(&ModelA{}).GetStructFields()
			} else {
				DB.NewScope(&ModelB{}).GetStructFields()
			}

			done.Done()
		}()

		start.Done()
	}

	done.Wait()
}

func TestModelStructEmbeddedHasMany(t *testing.T) {
	fields := DB.NewScope(&ResponseModel{}).GetStructFields()

	var childrenField *orm.StructField

	for i := 0; i < len(fields); i++ {
		field := fields[i]

		if field != nil && field.Name == "Children" {
			childrenField = field
		}
	}

	if childrenField == nil {
		t.Error("childrenField should not be nil")
		return
	}

	if childrenField.Relationship == nil {
		t.Error("childrenField.Relation should not be nil")
		return
	}

	expected := "has_many"
	actual := childrenField.Relationship.Kind

	if actual != expected {
		t.Errorf("childrenField.Relationship.Kind should be %v, but was %v", expected, actual)
	}
}
