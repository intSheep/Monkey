package parser

import (
	"testing"
)

func CheckErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %v errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error : %v", msg)
	}
	t.FailNow()
}
