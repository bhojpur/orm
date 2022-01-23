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
	"time"
)

func TestDelete(t *testing.T) {
	user1, user2 := User{Name: "delete1"}, User{Name: "delete2"}
	DB.Save(&user1)
	DB.Save(&user2)

	if err := DB.Delete(&user1).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	}

	if !DB.Where("name = ?", user1.Name).First(&User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}

	if DB.Where("name = ?", user2.Name).First(&User{}).RecordNotFound() {
		t.Errorf("Other users that not deleted should be found-able")
	}
}

func TestInlineDelete(t *testing.T) {
	user1, user2 := User{Name: "inline_delete1"}, User{Name: "inline_delete2"}
	DB.Save(&user1)
	DB.Save(&user2)

	if DB.Delete(&User{}, user1.Id).Error != nil {
		t.Errorf("No error should happen when delete a record")
	} else if !DB.Where("name = ?", user1.Name).First(&User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}

	if err := DB.Delete(&User{}, "name = ?", user2.Name).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	} else if !DB.Where("name = ?", user2.Name).First(&User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}
}

func TestSoftDelete(t *testing.T) {
	type User struct {
		Id        int64
		Name      string
		DeletedAt *time.Time
	}
	DB.AutoMigrate(&User{})

	user := User{Name: "soft_delete"}
	DB.Save(&user)
	DB.Delete(&user)

	if DB.First(&User{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	if err := DB.Unscoped().First(&User{}, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but err=%s", err)
	}

	DB.Unscoped().Delete(&user)
	if !DB.Unscoped().First(&User{}, "name = ?", user.Name).RecordNotFound() {
		t.Errorf("Can't find permanently deleted record")
	}
}

func TestSoftDeleteWithCustomizedDeletedAtColumnName(t *testing.T) {
	creditCard := CreditCard{Number: "411111111234567"}
	DB.Save(&creditCard)
	DB.Delete(&creditCard)

	if deletedAtField, ok := DB.NewScope(&CreditCard{}).FieldByName("DeletedAt"); !ok || deletedAtField.DBName != "deleted_time" {
		t.Errorf("CreditCard's DeletedAt's column name should be `deleted_time`")
	}

	if DB.First(&CreditCard{}, "number = ?", creditCard.Number).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	if err := DB.Unscoped().First(&CreditCard{}, "number = ?", creditCard.Number).Error; err != nil {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but err=%s", err)
	}

	DB.Unscoped().Delete(&creditCard)
	if !DB.Unscoped().First(&CreditCard{}, "number = ?", creditCard.Number).RecordNotFound() {
		t.Errorf("Can't find permanently deleted record")
	}
}
