import React, { useState, useEffect } from "react";

const App = () => {
    
    const [simulations, setSimulations] = useState([])
    useEffect(() => {
        fetch('http://localhost:8080/simulations')
            .then(response => response.json())
            .then(data => setSimulations(data.simulations))
    }, [])

    const stopSimulation = (simulation) => {
        console.log(simulation)
        fetch('http://localhost:8080/simulations/' + simulation.UserId + '/stop', { method: 'PUT' })
            .then(response => response.json())
            .then(data => setSimulations(data.simulations))
    }

    const resumeSimulation = (simulation) => {
        console.log(simulation)
        fetch('http://localhost:8080/simulations/' + simulation.UserId + '/resume', { method: 'PUT' })
            .then(response => response.json())
            .then(data => setSimulations(data.simulations))
    }

    const startNewSimulation = () => {
        fetch('http://localhost:8080/simulations', { method: 'POST' })
            .then(response => response.json())
            .then(data => setSimulations(data.simulations))
    }

    return (<div>
        <h1>SIMULATION</h1>
        <div>
        <button onClick={() => startNewSimulation()}>Start new user simulation</button>
        </div>
        <div>
                {simulations.map(simulation => (<div key={simulation.UserId}>
                    {simulation.UserId} - {simulation.Running ? "Running" : "Not running"} 

                    {simulation.Running ? 
                        <button onClick={() => stopSimulation(simulation)}>Stop</button> :
                        <button onClick={() => resumeSimulation(simulation)}>Resume</button>}
                </div>))}
        </div>
    </div>
    )
}

export default App;