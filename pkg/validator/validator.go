package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

type Validator struct {
	url.Values
	Errors errors
}

func New(form url.Values) *Validator {
	return &Validator{
		form,
		map[string][]string{},
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) MinLength(field string, n int) {
	val := v.Get(field)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) < n {
		v.Errors.Add(field, fmt.Sprintf("minimum number of characters is %d", n))
	}
}

func (v *Validator) MaxLength(field string, n int) {
	val := v.Get(field)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) > n {
		v.Errors.Add(field, fmt.Sprintf("maximum number of characters is %d", n))
	}
}

func (v *Validator) ValidEmail(field string) {
	val := v.Get(field)
	if val == "" {
		return
	}
	if !EmailRegex.MatchString(val) {
		v.Errors.Add(field, fmt.Sprintf("%s is not a valid email", val))
	}
}

func (v *Validator) Required(fields ...string) {
	for _, field := range fields {
		val := v.Get(field)
		if strings.TrimSpace(val) == "" {
			v.Errors.Add(field, "is required")
		}
	}
}

func (v *Validator) StrToInt(field string) int {
	val := v.Get(field)
	if val == "" {
		v.Errors.Add(field, "field is blank")
		return 0
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		v.Errors.Add(field, "field is not an integer")
		return 0
	}

	return i
}

func (v *Validator) ValidReaction(field string) {
	val := v.Get(field)
	if val != "like" && val != "dislike" {
		v.Errors.Add(field, "reaction is invalid")
	}
}
