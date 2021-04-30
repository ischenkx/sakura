package pubsub

type LocalStorage interface {
	Load(key interface{}) (value interface{}, ok bool)
	Store(key interface{}, value interface{})
	LoadOrStore(key interface{}, value interface{}) (actual interface{}, loaded bool)
	LoadAndDelete(key interface{}) (value interface{}, loaded bool)
	Delete(key interface{})
	Range(f func(key interface{}, value interface{}) bool)
}