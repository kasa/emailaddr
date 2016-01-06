package emailaddr

import (
	"bytes"
	"io"
	"regexp"

	"github.com/golang/glog"
)

const (
	initialize int = iota
	commentLocalBeg
	local              // 2
	localDot           //3
	localQuote         //4
	localDotQuoteStart //3
	localDotQuoteEnd   //3
	localEscape        //5
	commentLocalEnd
	at
	commentDomainBeg
	domain
	domainDot
	domainQuote
	commentDomainEnd
	end
	error
)

var hostChars = regexp.MustCompile("")

func IsValid(email string) bool {
	if len(email) == 0 {
		return false
	}
	if email[0] == ' ' || email[len(email)-1] == ' ' {
		return false
	}
	buf := bytes.NewBufferString(email)
	state := initialize
	i := 0
	for {
		r, _, err := buf.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			glog.Fatal(err)
		}

		// fmt.Println(string(r), state)

		switch state {
		case initialize:
			if r == '(' {
				state = commentLocalBeg
			} else if r == '"' {
				state = localQuote
			} else {
				buf.UnreadRune()
				state = local
			}
		case commentLocalBeg:
			if r == ')' {
				state = local
			}
		case local:
			if i > 64 {
				return false
			}
			switch r {
			case '(':
				state = commentLocalEnd
			case '@':
				state = at
			case '.':
				state = localDot
			case ' ', '"', ')', ',', ':', ';', '<', '>', '[', ']':
				state = error
			}
			i++
		case localDot:
			if r == '.' {
				state = error
			} else if r == '"' {
				state = localDotQuoteStart
			} else {
				buf.UnreadRune()
				state = local
			}
			i++
		case localQuote:
			if r == '"' {
				state = local
			} else if r == '\\' {
				state = localEscape
			}
		case localDotQuoteStart:
			if r == '"' {
				state = localDotQuoteEnd
			}
		case localDotQuoteEnd:
			if r == '.' {
				state = local
			} else {
				state = error
			}
		case localEscape:
			state = localQuote
		case at:
			i = 0
			if r == '(' {
				state = commentDomainBeg
			} else if r == ')' {
				state = domain
			} else {
				state = domain
			}
		case domain:
			if i > 253 {
				state = error
			} else if (r < 48 && r != '-' && r != '.') || (r > 71 && r < 65) || (r > 90 && r < 61) || r > 122 {
				state = error
			} else if r == '(' {
				state = commentDomainEnd
			} else if r == '@' {
				return false
			} else if r == '.' {
				state = domainDot
			}
		case domainDot:
			if r == '.' {
				state = error
			} else {
				buf.UnreadRune()
				state = domain
			}
		case commentDomainEnd:
			if r == ')' {
				state = end
			}
		case error:
			break
		}
	}

	return state == end || state == domain
}
