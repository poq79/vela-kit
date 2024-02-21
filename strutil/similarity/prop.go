package similarity

type Prop struct {
	v1       []rune
	v2       []rune
	max      int
	overlapN int
	overlap  map[rune]int
}

func NewProportion() *Prop {
	return &Prop{}
}

func (p *Prop) Compare(a, b string) float64 {
	v1 := []rune(a)
	v2 := []rune(b)
	v1n := len(v1)
	v2n := len(v2)
	if v1n == 0 || v2n == 0 {
		return 0
	}

	long := v1
	short := v2
	longN := v1n
	shortN := v2n

	if v1n < v2n {
		long = v2
		short = v1
		longN = v2n
		shortN = v1n
	}

	overlap := make(map[rune]bool, longN/2)
	overlapN := 0

	for i := 0; i < longN; i++ {
		overlap[long[i]] = true
	}

	for i := 0; i < shortN; i++ {
		if _, ok := overlap[short[i]]; ok {
			overlapN++
		}
	}

	return float64(overlapN) / float64(longN)
}
