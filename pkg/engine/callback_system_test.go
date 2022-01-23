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
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func equalFuncs(funcs []*func(s *Scope), fnames []string) bool {
	var names []string
	for _, f := range funcs {
		fnames := strings.Split(runtime.FuncForPC(reflect.ValueOf(*f).Pointer()).Name(), ".")
		names = append(names, fnames[len(fnames)-1])
	}
	return reflect.DeepEqual(names, fnames)
}

func create(s *Scope)        {}
func beforeCreate1(s *Scope) {}
func beforeCreate2(s *Scope) {}
func afterCreate1(s *Scope)  {}
func afterCreate2(s *Scope)  {}

func TestRegisterCallback(t *testing.T) {
	var callback = &Callback{logger: defaultLogger}

	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("before_create2", beforeCreate2)
	callback.Create().Register("create", create)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Register("after_create2", afterCreate2)

	if !equalFuncs(callback.creates, []string{"beforeCreate1", "beforeCreate2", "create", "afterCreate1", "afterCreate2"}) {
		t.Errorf("register callback")
	}
}

func TestRegisterCallbackWithOrder(t *testing.T) {
	var callback1 = &Callback{logger: defaultLogger}
	callback1.Create().Register("before_create1", beforeCreate1)
	callback1.Create().Register("create", create)
	callback1.Create().Register("after_create1", afterCreate1)
	callback1.Create().Before("after_create1").Register("after_create2", afterCreate2)
	if !equalFuncs(callback1.creates, []string{"beforeCreate1", "create", "afterCreate2", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}

	var callback2 = &Callback{logger: defaultLogger}

	callback2.Update().Register("create", create)
	callback2.Update().Before("create").Register("before_create1", beforeCreate1)
	callback2.Update().After("after_create2").Register("after_create1", afterCreate1)
	callback2.Update().Before("before_create1").Register("before_create2", beforeCreate2)
	callback2.Update().Register("after_create2", afterCreate2)

	if !equalFuncs(callback2.updates, []string{"beforeCreate2", "beforeCreate1", "create", "afterCreate2", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}
}

func TestRegisterCallbackWithComplexOrder(t *testing.T) {
	var callback1 = &Callback{logger: defaultLogger}

	callback1.Query().Before("after_create1").After("before_create1").Register("create", create)
	callback1.Query().Register("before_create1", beforeCreate1)
	callback1.Query().Register("after_create1", afterCreate1)

	if !equalFuncs(callback1.queries, []string{"beforeCreate1", "create", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}

	var callback2 = &Callback{logger: defaultLogger}

	callback2.Delete().Before("after_create1").After("before_create1").Register("create", create)
	callback2.Delete().Before("create").Register("before_create1", beforeCreate1)
	callback2.Delete().After("before_create1").Register("before_create2", beforeCreate2)
	callback2.Delete().Register("after_create1", afterCreate1)
	callback2.Delete().After("after_create1").Register("after_create2", afterCreate2)

	if !equalFuncs(callback2.deletes, []string{"beforeCreate1", "beforeCreate2", "create", "afterCreate1", "afterCreate2"}) {
		t.Errorf("register callback with order")
	}
}

func replaceCreate(s *Scope) {}

func TestReplaceCallback(t *testing.T) {
	var callback = &Callback{logger: defaultLogger}

	callback.Create().Before("after_create1").After("before_create1").Register("create", create)
	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Replace("create", replaceCreate)

	if !equalFuncs(callback.creates, []string{"beforeCreate1", "replaceCreate", "afterCreate1"}) {
		t.Errorf("replace callback")
	}
}

func TestRemoveCallback(t *testing.T) {
	var callback = &Callback{logger: defaultLogger}

	callback.Create().Before("after_create1").After("before_create1").Register("create", create)
	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Remove("create")

	if !equalFuncs(callback.creates, []string{"beforeCreate1", "afterCreate1"}) {
		t.Errorf("remove callback")
	}
}
