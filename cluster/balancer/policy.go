package balancer


type Policy interface {
	Next(AddressList) (Instance, error)
}
