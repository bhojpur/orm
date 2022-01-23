package postgres

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
	"database/sql/driver"

	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

type Hstore map[string]*string

// Value get value of Hstore
func (h Hstore) Value() (driver.Value, error) {
	hstore := hstore.Hstore{Map: map[string]sql.NullString{}}
	if len(h) == 0 {
		return nil, nil
	}

	for key, value := range h {
		var s sql.NullString
		if value != nil {
			s.String = *value
			s.Valid = true
		}
		hstore.Map[key] = s
	}
	return hstore.Value()
}

// Scan scan value into Hstore
func (h *Hstore) Scan(value interface{}) error {
	hstore := hstore.Hstore{}

	if err := hstore.Scan(value); err != nil {
		return err
	}

	if len(hstore.Map) == 0 {
		return nil
	}

	*h = Hstore{}
	for k := range hstore.Map {
		if hstore.Map[k].Valid {
			s := hstore.Map[k].String
			(*h)[k] = &s
		} else {
			(*h)[k] = nil
		}
	}

	return nil
}

// Jsonb Postgresql's JSONB data type
type Jsonb struct {
	json.RawMessage
}

// Value get value of Jsonb
func (j Jsonb) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan scan value into Jsonb
func (j *Jsonb) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}
