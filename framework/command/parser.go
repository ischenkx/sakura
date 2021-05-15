package command

import (
	"errors"
	"fmt"
	"strings"
)

const expressionPrefix = "notify:"

const nameCharset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@$_"

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}

	return s[:len(prefix)] == prefix
}

func IsExpression(s string) bool {
	return startsWith(s, expressionPrefix)
}

func getExpressionName(s string) (string, string, error) {
	if !IsExpression(s) {
		return "", "", errors.New("not an expression")
	}

	s = s[len(expressionPrefix):]

	nameLength := 0

	for nameLength < len(s) {
		if !strings.ContainsRune(nameCharset, rune(s[nameLength])) {
			break
		}
		nameLength++
	}

	if nameLength == 0 {
		return "", "", errors.New("no name provided")
	}

	return s[:nameLength], s[nameLength:], nil
}



type lexemInfo struct {
	typ, val string
}

func lexem(typ, val string) lexemInfo {
	return lexemInfo{typ, val}
}

func lexFlags(s string) ([]lexemInfo, error) {
	inQuotes := false
	curValue := ""

	var result []lexemInfo

	for i := 0; i < len(s); i++ {
		c := rune(s[i])

		if inQuotes {
			if c == '"' && rune(s[i-1]) != '\\' {
				inQuotes = false
				result = append(result, lexem("str", curValue))
				curValue = ""
				continue
			}
			curValue += string(c)
			continue
		}

		switch c {
		case ' ':
			if len(curValue) > 0 {
				result = append(result, lexem("lit", curValue))
				curValue = ""
			}
			continue
		case '=':
			if len(curValue) > 0 {
				result = append(result, lexem("lit", curValue))
				curValue = ""
			}
			result = append(result, lexem("eq", "="))
		case '"':
			inQuotes = true
		default:
			curValue += string(c)
		}
	}
	if inQuotes {
		return nil, errors.New("not closed quote")
	}
	if len(curValue) > 0 {
		result = append(result, lexem("lit", curValue))
		curValue = ""
	}
	return result, nil
}

// notify:event use_middleware=true name="Hello World"

func getFlags(s string) ([]Flag, error) {

	var flags []Flag

	lexems, err := lexFlags(s)

	if err != nil {
		return nil, err
	}

	currentFlag := Flag{}

	for i := 0; i < len(lexems); i++ {
		lex := lexems[i]

		switch lex.typ {
		case "str":
			if i > 0 {
				prev := lexems[i-1]
				if prev.typ == "eq" {
					currentFlag.Value = lex.val
					flags = append(flags, currentFlag)
					currentFlag = Flag{}
					continue
				}
			}
			flags = append(flags, Flag{Value: lex.val})
		case "lit":
			if i > 0 {
				prev := lexems[i-1]
				if prev.typ == "eq" {
					currentFlag.Value = lex.val
					flags = append(flags, currentFlag)
					currentFlag = Flag{}
					continue
				}
			}
			currentFlag.Name = lex.val
			if i == len(lexems) - 1 {
				flags = append(flags, currentFlag)
			}
		case "eq":
			if i == 0 {
				return nil, errors.New("= cannot be the first lexem in a flag line")
			}
			prev := lexems[i-1]
			if prev.typ != "lit" {
				return nil, errors.New("= can only go after a literal")
			}
			currentFlag.hasEq = true
		}
	}

	return flags, nil
}

func Parse(s string) (expr Command, err error) {
	s = strings.TrimSpace(s)

	name, otherData, err := getExpressionName(s)

	if err != nil {
		return expr, err
	}

	expr.Command = name

	flags, err := getFlags(otherData)

	if err != nil {
		return expr, err
	}
	for i := 0; i < len(flags); i++ {
		f := flags[i]
		if f.Name == "" {
			return expr, fmt.Errorf("can not find a name for a flag with value \"%s\"", f.Value)
		}
		if f.Value == "" && f.hasEq {
			return expr, errors.New("unexpected equality sign")
		}
	}
	expr.Flags = flags
	return expr, err
}

func IterText(txt string, f func(Command, error)) {
	comments := strings.Split(txt, "\n")
	for _, com := range comments {
		if !IsExpression(com) {
			continue
		}
		exp, err := Parse(com)
		f(exp, err)
	}
}