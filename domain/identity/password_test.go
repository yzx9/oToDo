package identity_test

import (
	"errors"
	"testing"

	"github.com/yzx9/otodo/domain/identity"
)

func TestPassword(t *testing.T) {
	t.Parallel()

	for _, test := range []string{"admin123", "passwd", "password", "abcdefghijklmnopqrst"} {
		pwd, err := identity.NewPassword(test)
		if err != nil {
			t.Errorf("Validate password with %v; got error(%s); expected: password", test, err.Error())
		} else if pwd.Empty() {
			t.Errorf("Validate password with %v; got empty password; expected: password", test)
		} else if !pwd.Equals(test) {
			t.Errorf("passowrd(%v).Equal(%v) == false; expected: true", test, test)
		}
	}

	if pwd, err := identity.NewPassword(""); err != nil || !pwd.Empty() {
		t.Errorf("Validate password with empty string; expected: empty password")
	}

	if _, err := identity.NewPassword("admin"); !errors.Is(err, identity.InvalidPassword) {
		t.Errorf("Validate password with short string; expected: password too short")
	}

	if _, err := identity.NewPassword("abcdefghijklmnopqrstu"); !errors.Is(err, identity.InvalidPassword) {
		t.Errorf("Validate password with long string; expected: password too long")
	}
}
