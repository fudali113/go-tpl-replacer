package main

import "testing"

func Test_loadKvString(t *testing.T) {
	testStrs := []string{
		"a=b",
	}
	argMap := map[string]interface{}{}
	for _, str := range testStrs {
		loadKvString(argMap, str)
	}
	_, ok := argMap["a"]
	if !ok {
		t.Error("fun loadKvString has error")
	}
}
