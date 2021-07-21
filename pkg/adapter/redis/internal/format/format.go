package format

import "fmt"

func Clients(ids []string) ([]string, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = Client(id)
	}

	return fmtIDS, true
}

func Users(ids []string) ([]string, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = User(id)
	}

	return fmtIDS, true
}

func Topics(ids []string) ([]string, bool) {

	if len(ids) == 0 {
		return nil, false
	}

	fmtIDS := make([]string, len(ids))

	for i, id := range ids {
		fmtIDS[i] = Topic(id)
	}

	return fmtIDS, true
}

func Client(id string) string {
	return fmt.Sprintf("c:%s", id)
}

func User(id string) string {
	return fmt.Sprintf("u:%s", id)
}

func Topic(id string) string {
	return fmt.Sprintf("t:%s", id)
}

