package model

type DataRow struct {
	RowIndex      int
	Project       string
	Type          string
	Key           string
	Resume        string
	Status        string
	Sprint        string
	ToStart       string
	InProgress    string
	InReview      string
	ToVerify      string
	Ready         string
	TimesInRework string
	StoryPoints   string
	Estimation    string
	SumEstimation string
}
