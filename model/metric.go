package model

type Metric struct {
	RowIndex        int
	Project         string
	Size            string
	NumHUS          int
	LeadTime        float64 //10-6
	PlanningTime    float64 //7-6
	CycleTime       float64 //9-7
	DevelopmentTime float64 //8-7
	VerifyingTime   float64 //10-9
	Deviation       float64
	ExecutionTime   float64
}
