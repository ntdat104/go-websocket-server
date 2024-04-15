package utils

type HashSet map[string]struct{}

func (set HashSet) Add(item string) {
	set[item] = struct{}{}
}

func (set HashSet) Remove(item string) {
	delete(set, item)
}

func (set HashSet) Contains(item string) bool {
	_, found := set[item]
	return found
}
