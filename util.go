package pbrtparser

import "strings"

func toTokens(s string) []string {
	// "NAME" "ParamName Type" [Values]
	// will be tokenize as
	// ", NAME, ", ", ParamName Type, ", [, Values, ]
	tokens := []string{}
	pushTok := func(t string) {
		t = strings.TrimSpace(t)
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	prev := ""
	for ptr := 0; ptr < len(s); ptr++ {
		ch := s[ptr : ptr+1]
		if ch == `[` || ch == `]` || ch == `"` {
			pushTok(prev)
			pushTok(ch)
			prev = ""
		} else {
			prev += ch
		}
	}
	pushTok(prev)
	return tokens
}

func assureForm(tokens []string, p0, p2 string) error {
	if tokens[0] != p0 || tokens[2] != p2 {
		return ErrClassCommandForm
	}
	return nil
}
