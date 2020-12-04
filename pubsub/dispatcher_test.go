package pubsub

import (
	"flag"
	"fmt"
	"github.com/RomanIschenko/notify/pubsub/internal/distributor"
	"github.com/RomanIschenko/notify/result"
	"github.com/google/uuid"
	"math/rand"
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	flag.Parse()

	c := distributor.New(16)

	topics := []string{}

	batch := distributor.Batch{}

	c.Aggregate(func(i int) result.Result {
		r := result.New()

		for i := 0; i < 1000; i++ {
			id := uuid.New().String()
			r.TopicsUp = append(r.TopicsUp, id)
			topics = append(topics, id)
		}

		return r
	})

	for _, t := range topics {
		if rand.Intn(10) > 5 {
			batch.Topics = append(batch.Topics, t)
		}

		if len(batch.Topics) > 4 {
			break
		}
	}

	fmt.Println("starting test", len(batch.Topics))

	start := time.Now()

	for i := 0; i < 100000; i++ {
		c.AggregateSharded(batch, func(shard int, data distributor.Batch) (r result.Result) {

			return
		})
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed, elapsed / 100000)

}
