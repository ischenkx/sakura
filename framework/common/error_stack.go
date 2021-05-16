package common

import "fmt"

func containsString(arr []string, s string) bool {
	for _, s1 := range arr {
		if s1 == s {
			return true
		}
	}
	return false
}

type MessageStack struct {
	errors, warnings []string
}

func (e *MessageStack) hasError(msg string) bool {
	return containsString(e.errors, msg)
}

func (e *MessageStack) hasWarning(msg string) bool {
	return containsString(e.warnings, msg)
}

func (e *MessageStack) Error(msg string) {
	if e.hasError(msg) {
		fmt.Println("contains!!!")
		return
	}
	e.errors = append(e.errors, msg)
}

func (e *MessageStack) Warning(msg string) {
	if e.hasWarning(msg) {return}

	e.warnings = append(e.warnings, msg)
}

func (e *MessageStack) Ok() bool {
	return len(e.errors) == 0
}

func (e *MessageStack) Errors() []string {
	return e.errors
}

func (e *MessageStack) Warnings() []string {
	return e.warnings
}
