package distributor

type Batch struct {
	Clients, Users, Topics []string
}

func (batch *Batch) Merge(b1 Batch) {
	batch.Topics = append(batch.Topics, b1.Topics...)
	batch.Clients = append(batch.Clients, b1.Clients...)
	batch.Users = append(batch.Users, b1.Users...)
}

func (batch *Batch) Reset() {
	batch.Topics = batch.Topics[:0]
	batch.Clients = batch.Clients[:0]
	batch.Users = batch.Users[:0]
}
