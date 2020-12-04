package balancer

import "time"

type Instance struct {
	ID, Address string
	time        time.Time
}