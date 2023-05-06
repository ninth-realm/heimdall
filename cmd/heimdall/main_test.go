package main

import "testing"

func TestStuff(t *testing.T) {
	if 1 == 2 {
		t.Fatal("Math is broken")
	}
}