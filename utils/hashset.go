package utils

type HashSet struct {
	set map[interface{}]bool
}

func NewHashSet() *HashSet {
	return &HashSet{set: make(map[interface{}]bool)}
}

func (s *HashSet) Add(key interface{}) {
	s.set[key] = true
}

func (s *HashSet) Remove(key interface{}) {
	delete(s.set, key)
}

func (s *HashSet) Contains(key interface{}) bool {
	_, found := s.set[key]
	return found
}

func (s *HashSet) Size() int {
	return len(s.set)
}

func (s *HashSet) Clear() {
	s.set = make(map[interface{}]bool)
}
