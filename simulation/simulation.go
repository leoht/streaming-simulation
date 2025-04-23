package simulation

type Simulation struct {
	userSimulations []*UserSimulation
}

var currentSimulation *Simulation

func StartSimulation() {
	currentSimulation = &Simulation{
		[]*UserSimulation{},
	}
}

func GetSimulations() []*UserSimulation {
	return currentSimulation.userSimulations
}
