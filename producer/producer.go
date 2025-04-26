package producer

import "leohetsch.com/simulation/simulation"

// A Client is anything with a Start() function
// accepting a channel to receive Events produced by the simulation
type Client interface {
	Start(producerInChannel chan simulation.Event)
}
