package emailaddr

import "testing"

func TestValid(t *testing.T) {
	var emails = []string{
		`(comment)prettyandsimple@example.com`,
		`prettyandsimple@example.com`,
		`very.common@example.com`,
		`disposable.style.email.with+symbol@example.com`,
		`other.email-with-dash@example.com`,
		`abc."defghi".xyz@example.com`,
		`"much.more unusual"@example.com`,
		`"very.unusual.@.unusual.com"@example.com`,
		`"very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com`,
		`admin@mailserver1`,
		"#!$%&'*+-/=?^_`{}|~@example.org",
		"\"()<>[]:,;@\\\"!#$%&'*+-/=?^_`{}| ~.a\"@example.org",
		"\" \"@example.org",
		`üñîçøðé@example.com`,
		`example@localhost`,
		`example@s.solutions`,
		`user@com`,
		`user@localserver`,
		`user@[IPv6:2001:db8::1]`,
	}

	for _, email := range emails {
		if !IsValid(email) {
			t.Fatalf("%s should be valid", email)
		}
	}
}

func TestInvalid(t *testing.T) {
	var emails = []string{
		`Abc.example.com`,
		`A@b@c@example.com`,
		`a"b(c)d,e:f;g<h>i[j\k]l@example.com`,
		`just"not"right@example.com`,
		`this is"not\allowed@example.com`,
		`this\ still\"not\\allowed@example.com`,
		`üñîçøðé@üñîçøðé.com`,
		`john..doe@example.com`,
		`john.doe@example..com`,
		`john doe@example.com`,
		` example@s.solutions`,
		`example@s.solutions `,
		"{example@s.solutions}",
		"",
	}

	for _, email := range emails {
		if IsValid(email) {
			t.Fatalf("%s should be invalid", email)
		}
	}
}
