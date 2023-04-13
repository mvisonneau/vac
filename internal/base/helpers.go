// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package base

import (
	"github.com/kr/text"
	"strings"
)

const (
	// maxLineLength is the maximum width of any line.
	maxLineLength int = 78
)

// SanitizePath removes any leading or trailing things from a "path".
func SanitizePath(s string) string {
	return EnsureNoTrailingSlash(EnsureNoLeadingSlash(s))
}

// EnsureTrailingSlash ensures the given string has a trailing slash.
func EnsureTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] != '/' {
		s = s + "/"
	}
	return s
}

// EnsureNoTrailingSlash ensures the given string does not have a trailing slash.
func EnsureNoTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

// EnsureNoLeadingSlash ensures the given string does not have a leading slash.
func EnsureNoLeadingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

// WrapAtLengthWithPadding wraps the given text at the maxLineLength, taking
// into account any provided left padding.
func WrapAtLengthWithPadding(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}

// WrapAtLength wraps the given text to maxLineLength.
func WrapAtLength(s string) string {
	return WrapAtLengthWithPadding(s, 0)
}
