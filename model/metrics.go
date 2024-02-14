package model

type Metrics struct {
	List []Metric
}

func (m *Metrics) Exists(project, size string) bool {
	for _, metrica := range m.List {
		if metrica.Project == project && metrica.Size == size {
			return true
		}
	}
	return false
}
func (m *Metrics) UpdateMetric(val Metric) bool {
	for i, metrica := range m.List {
		if metrica.Project == val.Project && metrica.Size == val.Size {

			m.List[i].NumHUS = metrica.NumHUS + 1
			m.List[i].LeadTime = (metrica.LeadTime + val.LeadTime)                      // float64(m.Projectes[i].NumHUS)
			m.List[i].PlanningTime = (metrica.PlanningTime + val.PlanningTime)          // float64(m.Projectes[i].NumHUS)
			m.List[i].CycleTime = (metrica.CycleTime + val.CycleTime)                   // float64(m.Projectes[i].NumHUS)
			m.List[i].DevelopmentTime = (metrica.DevelopmentTime + val.DevelopmentTime) // float64(m.Projectes[i].NumHUS)
			m.List[i].VerifyingTime = (metrica.VerifyingTime + val.VerifyingTime)       // float64(m.Projectes[i].NumHUS)

			return true // Return true to indicate that an update was made
		}
	}
	return false // Return false if no matching Metrica was found
}

func (m *Metrics) AddMetric(newMetrica Metric) {
	m.List = append(m.List, newMetrica)
}