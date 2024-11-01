package main

import (
	"fmt"
	"testing"
)

func TestParseJson(t *testing.T) {
	d, err := ParseJson(`(1+1)^2+1)`, true)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(d) // map[a:1 b:str c:[1 str] d:map[d1:1 d2:str] e:[] f:map[]]
}
