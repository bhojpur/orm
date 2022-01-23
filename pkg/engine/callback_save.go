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
	"strings"
)

func beginTransactionCallback(scope *Scope) {
	scope.Begin()
}

func commitOrRollbackTransactionCallback(scope *Scope) {
	scope.CommitOrRollback()
}

func saveAssociationCheck(scope *Scope, field *Field) (autoUpdate bool, autoCreate bool, saveReference bool, r *Relationship) {
	checkTruth := func(value interface{}) bool {
		if v, ok := value.(bool); ok && !v {
			return false
		}

		if v, ok := value.(string); ok {
			v = strings.ToLower(v)
			return v == "true"
		}

		return true
	}

	if scope.changeableField(field) && !field.IsBlank && !field.IsIgnored {
		if r = field.Relationship; r != nil {
			autoUpdate, autoCreate, saveReference = true, true, true

			if value, ok := scope.Get("orm:save_associations"); ok {
				autoUpdate = checkTruth(value)
				autoCreate = autoUpdate
				saveReference = autoUpdate
			} else if value, ok := field.TagSettingsGet("SAVE_ASSOCIATIONS"); ok {
				autoUpdate = checkTruth(value)
				autoCreate = autoUpdate
				saveReference = autoUpdate
			}

			if value, ok := scope.Get("orm:association_autoupdate"); ok {
				autoUpdate = checkTruth(value)
			} else if value, ok := field.TagSettingsGet("ASSOCIATION_AUTOUPDATE"); ok {
				autoUpdate = checkTruth(value)
			}

			if value, ok := scope.Get("orm:association_autocreate"); ok {
				autoCreate = checkTruth(value)
			} else if value, ok := field.TagSettingsGet("ASSOCIATION_AUTOCREATE"); ok {
				autoCreate = checkTruth(value)
			}

			if value, ok := scope.Get("orm:association_save_reference"); ok {
				saveReference = checkTruth(value)
			} else if value, ok := field.TagSettingsGet("ASSOCIATION_SAVE_REFERENCE"); ok {
				saveReference = checkTruth(value)
			}
		}
	}

	return
}

func saveBeforeAssociationsCallback(scope *Scope) {
	for _, field := range scope.Fields() {
		autoUpdate, autoCreate, saveReference, relationship := saveAssociationCheck(scope, field)

		if relationship != nil && relationship.Kind == "belongs_to" {
			fieldValue := field.Field.Addr().Interface()
			newScope := scope.New(fieldValue)

			if newScope.PrimaryKeyZero() {
				if autoCreate {
					scope.Err(scope.NewDB().Save(fieldValue).Error)
				}
			} else if autoUpdate {
				scope.Err(scope.NewDB().Save(fieldValue).Error)
			}

			if saveReference {
				if len(relationship.ForeignFieldNames) != 0 {
					// set value's foreign key
					for idx, fieldName := range relationship.ForeignFieldNames {
						associationForeignName := relationship.AssociationForeignDBNames[idx]
						if foreignField, ok := scope.New(fieldValue).FieldByName(associationForeignName); ok {
							scope.Err(scope.SetColumn(fieldName, foreignField.Field.Interface()))
						}
					}
				}
			}
		}
	}
}

func saveAfterAssociationsCallback(scope *Scope) {
	for _, field := range scope.Fields() {
		autoUpdate, autoCreate, saveReference, relationship := saveAssociationCheck(scope, field)

		if relationship != nil && (relationship.Kind == "has_one" || relationship.Kind == "has_many" || relationship.Kind == "many_to_many") {
			value := field.Field

			switch value.Kind() {
			case reflect.Slice:
				for i := 0; i < value.Len(); i++ {
					newDB := scope.NewDB()
					elem := value.Index(i).Addr().Interface()
					newScope := newDB.NewScope(elem)

					if saveReference {
						if relationship.JoinTableHandler == nil && len(relationship.ForeignFieldNames) != 0 {
							for idx, fieldName := range relationship.ForeignFieldNames {
								associationForeignName := relationship.AssociationForeignDBNames[idx]
								if f, ok := scope.FieldByName(associationForeignName); ok {
									scope.Err(newScope.SetColumn(fieldName, f.Field.Interface()))
								}
							}
						}

						if relationship.PolymorphicType != "" {
							scope.Err(newScope.SetColumn(relationship.PolymorphicType, relationship.PolymorphicValue))
						}
					}

					if newScope.PrimaryKeyZero() {
						if autoCreate {
							scope.Err(newDB.Save(elem).Error)
						}
					} else if autoUpdate {
						scope.Err(newDB.Save(elem).Error)
					}

					if !scope.New(newScope.Value).PrimaryKeyZero() && saveReference {
						if joinTableHandler := relationship.JoinTableHandler; joinTableHandler != nil {
							scope.Err(joinTableHandler.Add(joinTableHandler, newDB, scope.Value, newScope.Value))
						}
					}
				}
			default:
				elem := value.Addr().Interface()
				newScope := scope.New(elem)

				if saveReference {
					if len(relationship.ForeignFieldNames) != 0 {
						for idx, fieldName := range relationship.ForeignFieldNames {
							associationForeignName := relationship.AssociationForeignDBNames[idx]
							if f, ok := scope.FieldByName(associationForeignName); ok {
								scope.Err(newScope.SetColumn(fieldName, f.Field.Interface()))
							}
						}
					}

					if relationship.PolymorphicType != "" {
						scope.Err(newScope.SetColumn(relationship.PolymorphicType, relationship.PolymorphicValue))
					}
				}

				if newScope.PrimaryKeyZero() {
					if autoCreate {
						scope.Err(scope.NewDB().Save(elem).Error)
					}
				} else if autoUpdate {
					scope.Err(scope.NewDB().Save(elem).Error)
				}
			}
		}
	}
}
