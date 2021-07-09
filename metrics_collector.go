package swirl

import "sync"

type localMetrics struct {
	collector *metricsCollector
}

func (l localMetrics) Clients() IDList {
	return variadicIdList{
		count: func(list IDList) int {
			return l.collector.countClients()
		},
		array: func(list IDList) []string {
			return l.collector.getClients()
		},
	}
}

func (l localMetrics) Topics() IDList {
	return variadicIdList{
		count: func(list IDList) int {
			return l.collector.countTopics()
		},
		array: func(list IDList) []string {
			return l.collector.getTopics()
		},
	}
}

func (l localMetrics) Users() IDList {
	return variadicIdList{
		count: func(list IDList) int {
			return l.collector.countUsers()
		},
		array: func(list IDList) []string {
			return l.collector.getUsers()
		},
	}
}

type metricsCollector struct {
	mu sync.RWMutex
	appEvents AppEvents
	clients, users, topics []string
}

func (m *metricsCollector) countClients() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

func (m *metricsCollector) countUsers() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.users)
}

func (m *metricsCollector) countTopics() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.topics)
}

func (m *metricsCollector) getClients() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]string, len(m.clients))
	copy(arr, m.clients)
	return arr
}


func (m *metricsCollector) getTopics() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]string, len(m.topics))
	copy(arr, m.topics)
	return arr
}


func (m *metricsCollector) getUsers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]string, len(m.users))
	copy(arr, m.users)
	return arr
}

func (m *metricsCollector) close() {
	m.appEvents.Close()
}

func (m *metricsCollector) get() Metrics {
	return localMetrics{m}
}

func (m *metricsCollector) initBindings() {
	m.appEvents.OnChange(func(app *App, log ChangeLog) {
		m.mu.Lock()
		defer m.mu.Unlock()

		for _, clientUp := range log.ClientsUp {
			m.clients = append(m.clients, clientUp)
		}

		for _, topicUp := range log.TopicsUp {
			m.topics = append(m.topics, topicUp)
		}

		for _, userUp := range log.UsersUp {
			m.users = append(m.users, userUp)
		}

		for _, clientDown := range log.ClientsDown {
			m.clients = deleteFromStringArray(m.clients, clientDown)
		}

		for _, topicDown := range log.TopicsDown {
			m.topics = deleteFromStringArray(m.topics, topicDown)
		}

		for _, userDown := range log.UsersDown {
			m.users = deleteFromStringArray(m.users, userDown)
		}
	})
}

func deleteFromStringArray(arr []string, t string) []string {
	for i := 0; i < len(arr); i++ {
		s := arr[i]
		if s == t {
			arr[i] = arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			i--
		}
	}
	return arr
}

func newMetricsCollector(app *App) *metricsCollector {
	evs := app.Events(-1)
	c := &metricsCollector{
		mu:        sync.RWMutex{},
		appEvents: evs,
		clients:   nil,
		users:     nil,
		topics:    nil,
	}
	c.initBindings()
	return c
}