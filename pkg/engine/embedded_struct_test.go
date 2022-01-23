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

import "testing"

type BasePost struct {
	Id    int64
	Title string
	URL   string
}

type Author struct {
	ID    string
	Name  string
	Email string
}

type HNPost struct {
	BasePost
	Author  `orm:"embedded_prefix:user_"` // Embedded struct
	Upvotes int32
}

type EngadgetPost struct {
	BasePost BasePost `orm:"embedded"`
	Author   Author   `orm:"embedded;embedded_prefix:author_"` // Embedded struct
	ImageUrl string
}

func TestPrefixColumnNameForEmbeddedStruct(t *testing.T) {
	dialect := DB.NewScope(&EngadgetPost{}).Dialect()
	engadgetPostScope := DB.NewScope(&EngadgetPost{})
	if !dialect.HasColumn(engadgetPostScope.TableName(), "author_id") || !dialect.HasColumn(engadgetPostScope.TableName(), "author_name") || !dialect.HasColumn(engadgetPostScope.TableName(), "author_email") {
		t.Errorf("should has prefix for embedded columns")
	}

	if len(engadgetPostScope.PrimaryFields()) != 1 {
		t.Errorf("should have only one primary field with embedded struct, but got %v", len(engadgetPostScope.PrimaryFields()))
	}

	hnScope := DB.NewScope(&HNPost{})
	if !dialect.HasColumn(hnScope.TableName(), "user_id") || !dialect.HasColumn(hnScope.TableName(), "user_name") || !dialect.HasColumn(hnScope.TableName(), "user_email") {
		t.Errorf("should has prefix for embedded columns")
	}
}

func TestSaveAndQueryEmbeddedStruct(t *testing.T) {
	DB.Save(&HNPost{BasePost: BasePost{Title: "news"}})
	DB.Save(&HNPost{BasePost: BasePost{Title: "hn_news"}})
	var news HNPost
	if err := DB.First(&news, "title = ?", "hn_news").Error; err != nil {
		t.Errorf("no error should happen when query with embedded struct, but got %v", err)
	} else if news.Title != "hn_news" {
		t.Errorf("embedded struct's value should be scanned correctly")
	}

	DB.Save(&EngadgetPost{BasePost: BasePost{Title: "engadget_news"}})
	var egNews EngadgetPost
	if err := DB.First(&egNews, "title = ?", "engadget_news").Error; err != nil {
		t.Errorf("no error should happen when query with embedded struct, but got %v", err)
	} else if egNews.BasePost.Title != "engadget_news" {
		t.Errorf("embedded struct's value should be scanned correctly")
	}

	if DB.NewScope(&HNPost{}).PrimaryField() == nil {
		t.Errorf("primary key with embedded struct should works")
	}

	for _, field := range DB.NewScope(&HNPost{}).Fields() {
		if field.Name == "BasePost" {
			t.Errorf("scope Fields should not contain embedded struct")
		}
	}
}

func TestEmbeddedPointerTypeStruct(t *testing.T) {
	type HNPost struct {
		*BasePost
		Upvotes int32
	}

	DB.Create(&HNPost{BasePost: &BasePost{Title: "embedded_pointer_type"}})

	var hnPost HNPost
	if err := DB.First(&hnPost, "title = ?", "embedded_pointer_type").Error; err != nil {
		t.Errorf("No error should happen when find embedded pointer type, but got %v", err)
	}

	if hnPost.Title != "embedded_pointer_type" {
		t.Errorf("Should find correct value for embedded pointer type")
	}
}
