package redis

import "fmt"

type formatter struct {}

func (f formatter) Clients(ids []string) ([]string, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = f.Client(id)
	}

	return fmtIDS, true
}

func (f formatter) Users(ids []string) ([]string, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = f.User(id)
	}

	return fmtIDS, true
}

func (f formatter) Topics(ids []string) ([]string, bool) {

	if len(ids) == 0 {
		return nil, false
	}

	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = f.Topic(id)
	}

	return fmtIDS, true
}

func (formatter) Client(id string) string {
	return fmt.Sprintf("c:%s", id)
}

func (formatter) User(id string) string {
	return fmt.Sprintf("u:%s", id)
}

func (formatter) Topic(id string) string {
	return fmt.Sprintf("t:%s", id)
}