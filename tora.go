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
// Version: 1.1.1
// Descript: struct tag:tora
//			 struct help

package tora

import (
	"reflect"
	"fmt"
	"errors"
)

func Trans(dst interface{}, src interface{}) (err error) {

	errStr = []string{}

	dstValue := reflect.ValueOf(dst)
	dstMethod := dstValue.MethodByName("ToraMain")

	srcValue := reflect.ValueOf(src)
	srcMethod := srcValue.MethodByName("ToraMain")

	if dstValue.Kind() != reflect.Ptr || srcValue.Kind() != reflect.Ptr {
		return errors.New("[err] It is not a pointer to struct! ")
	}

	args := make([]reflect.Value, 0)
	var srcToraMain, dstToraMain = false, false

	if srcMethod.Kind() == reflect.Func {
		srcToraMain = srcMethod.Call(args)[0].Bool()
	}
	if dstMethod.Kind() == reflect.Func {
		dstToraMain = dstMethod.Call(args)[0].Bool()
	}

	if srcToraMain || dstToraMain == false && srcToraMain == false {
		return process(dst, src, dstValue, srcValue)
	} else {
		return process(src, dst, srcValue, dstValue)
	}
}

var masterTypeElem, slaveTypeElem reflect.Type
var masterValueElem, slaveValueElem reflect.Value

var errStr []string

// process
func process(slave, master interface{}, sValue, mValue reflect.Value) (err error) {

	masterTypeElem = reflect.TypeOf(master).Elem()
	masterValueElem = mValue.Elem()

	slaveTypeElem = reflect.TypeOf(slave).Elem()
	slaveValueElem = sValue.Elem()

	// must be struct
	if masterValueElem.Kind() != reflect.Struct || slaveValueElem.Kind() != reflect.Struct {
		return errors.New("[err] Pointer doesn't point to struct! ")
	}

	for i := 0; i < masterTypeElem.NumField(); i++ {
		// get tag
		tagName := masterTypeElem.Field(i).Tag.Get("tora")
		// if struct have tag
		if tagName != "" {

			core(i, tagName, true)

		} else {

			core(i, masterTypeElem.Field(i).Name, false)

		}
	}

	if len(errStr) > 0 {
		return errors.New(fmt.Sprint(errStr))
	}
	return nil
}

func core(index int, fieldNameStr string, tag bool) {

	slaveElemField, has := slaveTypeElem.FieldByName(fieldNameStr)
	if has && masterTypeElem.Field(index).Type == slaveElemField.Type {

		fieldName := slaveValueElem.FieldByName(fieldNameStr)
		if fieldName.IsValid() && fieldName.CanSet() {
			fieldName.Set(masterValueElem.Field(index))
		} else {
			errStr = append(errStr, fmt.Sprintf("[err] '%s' field is invalid or can't be set! \n", fieldNameStr))
		}

	} else {

		if tag {
			// tag
			errStr = append(errStr, fmt.Sprintf("[err] '%s' tag name don't exist or some of error! \n", fieldNameStr))
		} else {
			// field
			errStr = append(errStr, fmt.Sprintf("[err] '%s' field don't exist or some of error! \n", fieldNameStr))
		}

	}
}
