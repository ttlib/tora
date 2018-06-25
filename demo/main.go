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
// Version: 1.0.0

package main

import (
	"fmt"
	"duang.wiki/utils/tora"
)

type userDO struct {
	PkId    int32 `tora:"Id" json:"Id"` // tag：tora
	Name    string
	Age     int
	Gender  string
	NetWork []string
}

// Implementing this interface is defined as the src
func (u *userDO) ToraMain() bool {
	return true
}

type userPB struct {
	Id      int32
	Name    string
	Age     string
	Gender  string
	NetWork []string
}

func main() {
	src := &userDO{10086, "hydden", 22, "male", []string{"wo", "ni"}}
	dst := &userPB{Name: "wihhy"}

	err := tora.Trans(dst, src)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("-----Modify src.name-----")
	src.Name = "Tom"
	fmt.Println(src)
	fmt.Println(dst)

	fmt.Println("-----Modify dst.Network-----")
	dst.NetWork = []string{"192.168.101.101"}
	fmt.Println(src)
	fmt.Println(dst)
}
