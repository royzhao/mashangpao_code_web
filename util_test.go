package main

import (
	"fmt"
	"testing"
)

func Test_setvalue(t *testing.T) {
	data := UserInfo{
		Id:      1,
		UserId:  1,
		Avatar:  "ghgh",
		Discrip: "fgtftft",
	}
	err := SetValue("uuuuu", data)
	if err != nil {
		t.Fatal(err)
	}
	status, data2 := GetValue("uuuuu")
	if status == 1 {
		t.Fatal(err)
	}
	fmt.Println(data2)

	status, _ = GetValue("uuuhhu")
	if status == 5 {
		t.Fatal(err)
	}
	// fmt.Println(data)
}
