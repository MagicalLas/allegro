package inmemory

var mainHashMap map[string]string

func init() {
	mainHashMap = make(map[string]string, 1000)
}

func Set(key, value string) {
	mainHashMap[key] = value
}

func Get(key string) (string, bool) {
	value, ok := mainHashMap[key]
	return value, ok
}
