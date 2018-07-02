// Copyright 2014 duang.wiki Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Tong Zicheng
// Date: 2018/6/25 11:43
// Descript: Model transformation help tools destination to source.
// 			 After this function, data is independent.
//           The destination's data is owerridden with default values.
//           Only support struct to struct or slice to slice.
//			 Struct's tag is `tora: ""`

package tora

import (
	"reflect"
	"fmt"
	"errors"
	"log"
)

// Print trans log
const (
	VERSION   = "1.4.0"
	TAG_NAME  = "tora"
	MAIN_FUNC = "ToraMain"

	LOG_FLAG = true
)

// Read the tag or fields for the struct that implements this interface.
type SelectTag interface {
	ToraMain() bool
}

// src -> dst Read src's tag or field by default. After this function, data is independent.
// The dst's data is owerridden with default values.
func Trans(dst interface{}, src interface{}) (err error) {

	// must be reflect.Ptr
	dstValue := reflect.ValueOf(dst) // ptr
	srcValue := reflect.ValueOf(src) // ptr

	if dstValue.Kind() != reflect.Ptr || srcValue.Kind() != reflect.Ptr {
		return errors.New("[err] It is not a pointer to struct! ")
	}

	dstValueElem := dstValue.Elem() // []*slice or struct
	srcValueElem := srcValue.Elem()

	dstType := reflect.TypeOf(dst) // ptr
	srcType := reflect.TypeOf(src)

	dstTypeElem := dstType.Elem() // []*slice or struct
	srcTypeElem := srcType.Elem()

	// Is the slice.
	if dstTypeElem.Kind() == reflect.Slice && srcTypeElem.Kind() == reflect.Slice {

		dstElemTypePtr := dstTypeElem.Elem() // ptr of []*slice[i]
		srcElemTypePtr := srcTypeElem.Elem()

		// slice item is ptr
		if dstElemTypePtr.Kind() == reflect.Ptr && srcElemTypePtr.Kind() == reflect.Ptr {

			for i := 0; i < srcValueElem.Len(); i++ {
				// Create a new ptr of dst's element.
				_dstValuePtr := reflect.New(dstElemTypePtr.Elem())
				_srcValuePtr := srcValueElem.Index(i)

				err := process(_dstValuePtr, _srcValuePtr, dstElemTypePtr, srcElemTypePtr)
				if err != nil {
					return err
				}

				dstValueElem.Set(reflect.Append(dstValueElem, _dstValuePtr))
			}

		}

		// slice item is struct
		if dstElemTypePtr.Kind() == reflect.Struct && srcElemTypePtr.Kind() == reflect.Struct {

			for i := 0; i < srcValueElem.Len(); i++ {
				// struct
				_dstValuePtr := reflect.New(dstElemTypePtr)
				_srcValuePtr := srcValueElem.Index(i).Addr()

				err := process(_dstValuePtr, _srcValuePtr, dstTypeElem, srcTypeElem)
				if err != nil {
					return err
				}

				dstValueElem.Set(reflect.Append(dstValueElem, _dstValuePtr.Elem()))
			}

		}
		// Is the struct.
	} else if dstTypeElem.Kind() == reflect.Struct && srcTypeElem.Kind() == reflect.Struct {

		err := process(dstValue, srcValue, dstType, srcType)
		if err != nil {
			return err
		}

		// Other Type
	} else {
		return errors.New("[err] There's no right function! ")
	}

	return nil

}

// A single struct transformation process.
func process(dstValue, srcValue reflect.Value, dstType, srcType reflect.Type) error {

	dstMethod := dstValue.MethodByName(MAIN_FUNC)
	srcMethod := srcValue.MethodByName(MAIN_FUNC)

	args := make([]reflect.Value, 0)
	var srcToraMainRes, dstToraMainRes = false, false

	// Determine if the struct has MAIN_FUNC and get the return value.
	if srcMethod.Kind() == reflect.Func {
		srcToraMainRes = srcMethod.Call(args)[0].Bool()
	}
	if dstMethod.Kind() == reflect.Func {
		dstToraMainRes = dstMethod.Call(args)[0].Bool()
	}

	if srcToraMainRes || dstToraMainRes == false && srcToraMainRes == false {
		// dst <- src  read SRC tag | default
		return parse(dstValue, srcValue, dstType, srcType, false)
	} else {
		// dst <- src read dst tag
		return parse(dstValue, srcValue, dstType, srcType, true)
	}
}

// Parse tag
func parse(dstValue, srcValue reflect.Value, dstType, srcType reflect.Type, dstTag bool) (err error) {

	srcTypeElem := srcType.Elem()
	srcValueElem := srcValue.Elem()

	dstTypeElem := dstType.Elem()
	dstValueElem := dstValue.Elem()

	// dst and src must be struct
	if srcValueElem.Kind() != reflect.Struct || dstValueElem.Kind() != reflect.Struct {
		return errors.New("[err] Pointer must be struct! ")
	}

	srcElemLen, dstElemLen := srcTypeElem.NumField(), dstTypeElem.NumField()

	var tagName = ""
	var dstTagNames map[string]string

	// Extract all dst's tags and value.
	if dstTag {
		dstTagNames = make(map[string]string)
		for j := 0; j < dstElemLen; j++ {

			tagName = dstTypeElem.Field(j).Tag.Get(TAG_NAME)

			switch tagName {
			case "":
			case "-":
				break
			default:
				dstTagNames[tagName] = dstTypeElem.Field(j).Name
			}

		}
	}

	// Iterate through all of the src's field.
	for i := 0; i < srcElemLen; i++ {

		if !dstTag {
			// dst <- src  READ SRC tag
			tagName = srcTypeElem.Field(i).Tag.Get(TAG_NAME)

			switch tagName {
			case "":
				core(dstValueElem, srcValueElem, dstTypeElem, srcTypeElem,
					i, srcTypeElem.Field(i).Name, 2)
				break
			case "-":
				break
			default:
				// If tag has value use it.
				core(dstValueElem, srcValueElem, dstTypeElem, srcTypeElem,
					i, tagName, 1)
			}

		} else {
			// dst <- src READ dst tag
			// If dst's tag has value, find the src's fieldname with it's value else use the value of dst's fieldname.
			if dstTagNames[srcTypeElem.Field(i).Name] != "" {

				core(dstValueElem, srcValueElem, dstTypeElem, srcTypeElem,
					i, dstTagNames[srcTypeElem.Field(i).Name], 3)

			} else {

				core(dstValueElem, srcValueElem, dstTypeElem, srcTypeElem,
					i, srcTypeElem.Field(i).Name, 3)

			}
		}
	}

	return nil
}

// Trans field
func core(dstValueElem, srcValueElem reflect.Value, dstTypeElem, srcTypeElem reflect.Type,
	srcIndex int, dstFieldNameStr string, tag int) {

	// model name
	dstNameStr := dstTypeElem.Name()
	srcNameStr := srcTypeElem.Name()

	// err message
	var errStr []string

	dstElemField, has := dstTypeElem.FieldByName(dstFieldNameStr) // Get Field
	if has && srcTypeElem.Field(srcIndex).Type == dstElemField.Type {

		dstFieldName := dstValueElem.FieldByName(dstFieldNameStr)

		// Field is valid and can be set
		if dstFieldName.IsValid() && dstFieldName.CanSet() {
			dstFieldValue := srcValueElem.Field(srcIndex)
			dstFieldName.Set(dstFieldValue)
		} else {
			errStr = append(errStr, fmt.Sprintf("[warn] '%s-%s' field is invalid or can't be set!", dstNameStr, dstFieldNameStr))
		}

	} else {

		// Wrong type，1 - src tag wrong, 2 - src filed name wrong, 3 - dst tag wrong
		switch tag {
		case 1:
			errStr = append(errStr, fmt.Sprintf("[warn] '%s-%s' tag name don't exist or have some error!", srcTypeElem, dstFieldNameStr))
			break
		case 2:
			errStr = append(errStr, fmt.Sprintf("[warn] '%s-%s' field don't exist or have some error!", srcNameStr, dstFieldNameStr))
			break
		case 3:
			errStr = append(errStr, fmt.Sprintf("[warn] (dst) '%s-%s' fields or tags hava some error!", dstNameStr, dstFieldNameStr))
			break
		}

	}

	if len(errStr) > 0 && LOG_FLAG {
		// return errors.New(fmt.Sprint(errStr))
		log.Println(errStr)
	}
}
