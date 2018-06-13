/*
 * Inflector Pkg (Go)
 *
 * Copyright (c) 2013 Ivan Torres
 * Released under the MIT license
 * https://github.com/mexpolk/inflector/blob/master/LICENSE
 *
 */

package strings

import (
	"regexp"
	"strings"
)

// Prepares strings by splitting by caps, spaces, dashes, and underscore
func split(str string) (words []string) {
	repl := strings.NewReplacer("-", " ", "_", " ")

	rex1 := regexp.MustCompile("([A-Z])")
	rex2 := regexp.MustCompile("(\\w+)")

	str = trim(str)

	// Convert dash and underscore to spaces
	str = repl.Replace(str)

	// Split when uppercase is found (needed for Snake)
	str = rex1.ReplaceAllString(str, " $1")

	// Get the final list of words
	words = rex2.FindAllString(str, -1)

	return
}

// Inflects to camelCase
func ToCamel(str string) (out string) {
	words := split(str)
	out = words[0]

	for _, w := range split(str)[1:] {
		out += ToUpper(w[:1]) + w[1:]
	}

	return
}

func ToLowerCamel(str string) (out string) {
	return ToCamel(str)
}

func ToUpperCamel(str string) (out string) {
	return ToPascal(str)
}

// Inflects to kebab-case
func ToDash(str string) (out string) {
	words := split(str)

	for i, w := range words {
		words[i] = ToLower(w)
	}

	out = strings.Join(words, "-")
	return
}

// Inflects to PascalCase
func ToPascal(str string) (out string) {
	out = ""

	for _, w := range split(str) {
		out += ToUpper(w[:1]) + w[1:]
	}

	return
}

// Inflects to snake_case
func ToSnake(str string) (out string) {
	words := split(str)

	for i, w := range words {
		words[i] = ToLower(w)
	}

	out = strings.Join(words, "_")
	return
}

// Inflects to snake_case
func ToTitle(str string) (out string) {

	out = ""
	words := split(str)

	for i, w := range words {
		words[i] = ToUpper(w[:1]) + w[1:]
	}

	out = strings.Join(words, " ")
	return
}

// Inflects to snake_case
func ToHeader(str string) (out string) {

	out = ""
	words := split(str)

	for i, w := range words {
		words[i] = ToUpper(w[:1]) + w[1:]
	}

	out = strings.Join(words, "-")
	return
}

// Removes leading whitespaces
func trim(str string) string {
	return strings.Trim(str, " ")
}

// Shortcut to strings.ToUpper()
func ToUpper(str string) string {
	return strings.ToUpper(trim(str))
}

// Shortcut to strings.ToLower()
func ToLower(str string) string {
	return strings.ToLower(trim(str))
}
