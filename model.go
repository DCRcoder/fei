package fei

type ExplainModel struct {
	ID           int64
	SelectType   string
	Table        string
	Partitions   *string
	Type         string
	PossibleKeys *string
	Key          *string
	KeyLen       *string
	Ref          *string
	Rows         int64
	Filtered     float64
	Extra        *string
}
