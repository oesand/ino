package validate_test

import (
	"regexp"
	"testing"

	"github.com/oesand/mo/validate"
)

func TestRegex(t *testing.T) {
	v := validate.Regex(regexp.MustCompile(`^abc$`))

	if res := v.Validate("abcd"); res.IsValid() {
		t.Fatalf("expected invalid for \"abcd\"")
	} else if err := res.Error(); err != "mismatch expected pattern" {
		t.Fatalf("unexpected error %s", err)
	}

	if res := v.Validate("abc"); !res.IsValid() {
		t.Fatalf("expected valid for \"abc\", got %v", res)
	}
}

func TestRegexWithMessage(t *testing.T) {
	v := validate.Regex(regexp.MustCompile(`^abc$`), validate.WithMessage("new message Regex"))

	if res := v.Validate("abcd"); res.IsValid() {
		t.Fatalf("expected invalid for \"abcd\"")
	} else if err := res.Error(); err != "new message Regex" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestRunesExactly(t *testing.T) {
	v := validate.RunesExactly(3)

	if res := v.Validate("ab"); res.IsValid() {
		t.Fatalf("expected invalid for \"ab\"")
	} else if err := res.Error(); err != "must have exactly 3 characters" {
		t.Fatalf("unexpected error %s", err)
	}

	if res := v.Validate("abc"); !res.IsValid() {
		t.Fatalf("expected valid for \"abc\", got %v", res)
	}
}

func TestRunesExactlyWithMessage(t *testing.T) {
	v := validate.RunesExactly(3, validate.WithMessage("new message RunesExactly"))

	if res := v.Validate("ab"); res.IsValid() {
		t.Fatalf("expected invalid for \"ab\"")
	} else if err := res.Error(); err != "new message RunesExactly" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestMinRunes(t *testing.T) {
	v := validate.MinRunes(2)

	if res := v.Validate("a"); res.IsValid() {
		t.Fatalf("expected invalid for \"a\"")
	} else if err := res.Error(); err != "must have at least 2 characters" {
		t.Fatalf("unexpected error %s", err)
	}

	if res := v.Validate("ab"); !res.IsValid() {
		t.Fatalf("expected valid for \"ab\", got %v", res)
	}
}

func TestMinRunesWithMessage(t *testing.T) {
	v := validate.MinRunes(2, validate.WithMessage("new message MinRunes"))

	if res := v.Validate("a"); res.IsValid() {
		t.Fatalf("expected invalid for \"a\"")
	} else if err := res.Error(); err != "new message MinRunes" {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestMaxRunes(t *testing.T) {
	v := validate.MaxRunes(2)

	if res := v.Validate("abc"); res.IsValid() {
		t.Fatalf("expected invalid for \"abc\"")
	} else if err := res.Error(); err != "must have at most 2 characters" {
		t.Fatalf("unexpected error %s", err)
	}

	if res := v.Validate("ab"); !res.IsValid() {
		t.Fatalf("expected valid for \"ab\", got %v", res)
	}
}

func TestMaxRunesWithMessage(t *testing.T) {
	v := validate.MaxRunes(2, validate.WithMessage("new message MaxRunes"))

	if res := v.Validate("abc"); res.IsValid() {
		t.Fatalf("expected invalid for \"abc\"")
	} else if err := res.Error(); err != "new message MaxRunes" {
		t.Fatalf("unexpected error %s", err)
	}
}
