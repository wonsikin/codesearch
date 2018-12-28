package src

// Records records
type Records []*Record

func (r Records) Len() int {
	return len(r)
}
func (r Records) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r Records) Less(i, j int) bool {
	return r[i].Count < r[j].Count
}
