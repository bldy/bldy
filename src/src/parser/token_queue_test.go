package parser

import "testing"

func TestTokenQueue(t *testing.T) {
	var q tq
	if q.pop() != nil {
		t.Fail()
	}
}
