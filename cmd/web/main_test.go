package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run() // なんかサラッとmain内の関数であるrunを使えているけど
	// privateではない状態だからいいみたいです
	if err != nil {
		t.Error("failed run()")
	}
}