package set

type Set[E comparable] map[E]struct{}

func New[E comparable]() Set[E] {
	return make(Set[E])
}

func (s Set[E]) Copy() Set[E] {
	c := New[E]()
	for e := range s {
		c.Add(e)
	}
	return c
}

func (s Set[E]) IsEmpty() bool {
	return len(s) == 0
}

func (s Set[E]) Len() int {
	return len(s)
}

func (s Set[E]) Contains(e E) bool {
	_, ok := s[e]
	return ok
}

func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
}

func (s Set[E]) Remove(e E) {
	delete(s, e)
}

func (s Set[E]) Union(t Set[E]) Set[E] {
	u := s.Copy()
	for e := range t {
		u.Add(e)
	}
	return u
}

func (s Set[E]) Intersect(t Set[E]) Set[E] {
	var needle, haystack Set[E]

	if s.Len() < t.Len() {
		needle, haystack = s, t
	} else {
		needle, haystack = t, s
	}

	i := New[E]()
	for e := range needle {
		if haystack.Contains(e) {
			i.Add(e)
		}
	}
	return i
}
