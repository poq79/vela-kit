package runtime

type Record struct {
	ID    int64   `json:"id" storm:"id"`
	Value float64 `json:"value"`
}
