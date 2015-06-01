package main

import (
	"fmt"
	"testing"
)

func Test_getUserinfo(t *testing.T) {
	var user UserInfo
	data, err := user.getInfo(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
	data2, err := user.getInfoFilter(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data2)
}
