package storages

func NewMemory[K comparable, V any]() map[K]V {
	return make(map[K]V)
}
