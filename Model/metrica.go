package model

type Metrica struct {
	Projecte        string
	Mida            string
	NumHUS          int
	LeadTime        float64 //10-6
	PlanningTime    float64 //7-6
	CycleTime       float64 //9-7
	DevelopmentTime float64 //8-7
	VerifyingTime   float64 //10-9
}

type Metriques struct {
	Projectes []Metrica
}

func (m *Metriques) Exists(projecte, mida string) bool {
	for _, metrica := range m.Projectes {
		if metrica.Projecte == projecte && metrica.Mida == mida {
			return true
		}
	}
	return false
}
func (m *Metriques) Update(val Metrica) bool {
	for i, metrica := range m.Projectes {
		if metrica.Projecte == val.Projecte && metrica.Mida == val.Mida {
			//fmt.Println("p- ",m.Projectes[i])
			//fmt.Println(metrica)
			// Match found, update the fields
			m.Projectes[i].NumHUS = metrica.NumHUS + 1
			m.Projectes[i].LeadTime = (metrica.LeadTime + val.LeadTime)                      // float64(m.Projectes[i].NumHUS)
			m.Projectes[i].PlanningTime = (metrica.PlanningTime + val.PlanningTime)          // float64(m.Projectes[i].NumHUS)
			m.Projectes[i].CycleTime = (metrica.CycleTime + val.CycleTime)                   // float64(m.Projectes[i].NumHUS)
			m.Projectes[i].DevelopmentTime = (metrica.DevelopmentTime + val.DevelopmentTime) // float64(m.Projectes[i].NumHUS)
			m.Projectes[i].VerifyingTime = (metrica.VerifyingTime + val.VerifyingTime)       // float64(m.Projectes[i].NumHUS)

			return true // Return true to indicate that an update was made
		}
	}
	return false // Return false if no matching Metrica was found
}

func (m *Metriques) AddMetrica(newMetrica Metrica) {
	m.Projectes = append(m.Projectes, newMetrica)
}