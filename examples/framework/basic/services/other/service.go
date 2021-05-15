package other

type Service struct {
	Data string
	// notify:inject label="dataService"
	DataService interface{}
}
