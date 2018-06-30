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
// Date: 2018/6/25 11:44

package main

import (
	"fmt"
	"github.com/ttlib/tora"
	"reflect"
)

type userDO struct {
	PkId    int32  `tora:"Id" json:"Id"` // tag：tora 定义 dst 的属性名
	Name    string `tora:"-"`
	Age     string
	Gender  bool
	NetWork []string
}

// dst 实现此接口，读取dst的 tag
func (u *userDO) ToraMain() bool {
	return true
}

type userPB struct {
	Id      int32
	Name    string
	Age     int
	Gender  string
	NetWork []string
}

func main() {
	//src4 := []*userPB{
	//	&userPB{01, "WHILE", 31, "男", []string{"a", "j"}},
	//	&userPB{02, "kaite", 33, "女", []string{"c", "f"}},
	//}
	//
	//src5 := []userPB{
	//	{01, "WHILE", 31, "男", []string{"a", "j"}},
	//	{02, "kaite", 33, "女", []string{"c", "f"}},
	//}

	demo1()
	demo2()
	demo3()
	demo4()
	//test(&src4 ,&src5)
}

func demo1() {
	// Demo1
	src := &userPB{10086, "颜如玉", 22, "女", []string{"wo", "ni"}}
	dst := &userDO{Name: "辗迟"}

	err := tora.Trans(dst, src)
	if err != nil {
		// 如果err有值，表示部分属性没有转换成功，类型错误，或者是tag错误
		fmt.Println(err)
	}

	fmt.Println("-----修改src.name-----")
	src.Name = "东方悦"
	fmt.Println(src)
	fmt.Println(dst)

	fmt.Println("-----修改dst.Network-----")
	dst.NetWork = []string{"192.168.101.101"}

	fmt.Println(src)
	fmt.Println(dst)
}

func demo2() {
	// Demo2
	src2 := []userPB{
		{01, "WHILE", 31, "男", []string{"a", "j"}},
		{02, "kaite", 33, "女", []string{"c", "f"}},
	}

	var dst2 []userDO

	err := tora.Trans(&dst2, &src2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(src2)
	fmt.Println(dst2)
}

func demo3() {
	// Demo3
	src3 := userPB{}
	dst3 := userDO{}

	err := tora.Trans(&dst3, &src3)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(src3)
	fmt.Println(dst3)
}

func demo4() {
	src4 := []*userPB{
		&userPB{01, "WHILE", 31, "男", []string{"a", "j"}},
		&userPB{02, "kaite", 33, "女", []string{"c", "f"}},
	}
	dst4 := []*userDO{}

	err := tora.Trans(&dst4,&src4)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range src4 {
		fmt.Println(v)
	}

	for _, v := range dst4 {
		fmt.Println(v)
	}
}

func test(src4,src5 interface{})  {

	//dst4 := []*userPB{}

	src4Type := reflect.TypeOf(src4).Elem().Elem().Elem()
	fmt.Println(reflect.TypeOf(src4).Elem().Kind()) //slice
	fmt.Println(reflect.TypeOf(src4).Elem().Elem().Kind()) //ptr

	fmt.Println(reflect.TypeOf(src5).Elem().Kind()) //slice
	fmt.Println(reflect.TypeOf(src5).Elem().Elem().Kind())//struct

	fmt.Println(src4Type)
	src4Obj := reflect.New(src4Type)
	fmt.Println(src4Obj)

	src4Value := reflect.ValueOf(src4).Elem()
	fmt.Println(src4Value)
	src4Value.Set(reflect.Append(src4Value,src4Obj))

	fmt.Println(src4)

	/*
	src4Value := reflect.ValueOf(src4)
	fmt.Println(src4Value)
	//dst4Value := reflect.ValueOf(dst4)
	//fmt.Println(dst4Value)
	indexData := src4Value.Index(0)
	fmt.Println(indexData)

	//value := reflect.Append(dst4Value, indexData)
	//value = reflect.Append(value, indexData)
	//value = reflect.Append(value, indexData)

	//indirect := reflect.Indirect(dst4Value)
	//slice := reflect.Append(indirect, indexData)

	dst4ValueP := reflect.ValueOf(&dst4).Elem()
	dst4ValueP.Set(reflect.Append(dst4ValueP,indexData))
	dst4ValueP.Set(reflect.Append(dst4ValueP,indexData))

	//
	fmt.Println(dst4)
	fmt.Println(src4)
	//dst4Value.SetPointer(dst4Value)
	*/
}

// Console print
/*
 *	-----修改src.name-----
 *	&{10086 东方悦 22 女 [wo ni]}
 *	&{10086 颜如玉  false [wo ni]}
 *	-----修改dst.Network-----
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Age' fields or tags hava some error!]
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Gender' fields or tags hava some error!]
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Age' fields or tags hava some error!]
 *	&{10086 东方悦 22 女 [wo ni]}
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Gender' fields or tags hava some error!]
 *	&{10086 颜如玉  false [192.168.101.101]}
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Age' fields or tags hava some error!]
 *	[{1 WHILE 31 男 [a j]} {2 kaite 33 女 [c f]}]
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Gender' fields or tags hava some error!]
 *	[{1 WHILE  false [a j]} {2 kaite  false [c f]}]
 *	{0  0  []}
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Age' fields or tags hava some error!]
 *	{0   false []}
 *	2018/06/29 10:41:47 [[warn] (dst) 'userDO-Gender' fields or tags hava some error!]
 */
