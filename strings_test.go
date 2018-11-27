package main

import (
	"testing"
)

//func Test_ (t *testing.T) {
//        t.Log()
//        t.Error()
//}

func Test_ComputeLine(t *testing.T) {
	lines := ComputeLine("./host1")
	if lines == 2 {
		t.Log("ComputeLine test passed")
	} else {
		t.Error("ComputeLine test failed: ", lines)
	}
}
