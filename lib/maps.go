package lib

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetMapKeys[K comparable, V any](m map[K]V, keys []K) []V {
	values := make([]V, 0, len(keys))
	for _, key := range keys {
		if v, found := m[key]; found {
			values = append(values, v)
		}
	}
	return values
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func MapEntries[K comparable, V any](m map[K]V) []Entry[K, V] {
	var rv []Entry[K, V]
	for k, v := range m {
		rv = append(rv, Entry[K, V]{Key: k, Value: v})
	}
	return rv
}
