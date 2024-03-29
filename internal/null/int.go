// Copyright 2023 GitBundle Inc. All rights reserved.
// Copyright (c) 2014, Greg Roseberry
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int is an nullable int64. It does not consider zero
// values to be null. It will decode to null, not zero,
// if null.
type Int struct {
	sql.NullInt64
}

// UnmarshalJSON implements json.Unmarshaler. It supports
// number and null input. 0 will not be considered a null
// Int. It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.Valid = false
			return nil
		}
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// IsZero returns true for invalid Ints, for future omitempty
// support (Go 1.4?). A non-null Int with a 0 value will not
// be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}
