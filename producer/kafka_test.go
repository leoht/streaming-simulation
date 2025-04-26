package producer

import (
	"fmt"
	"sync"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/magiconair/properties/assert"
	"leohetsch.com/simulation/simulation"
)

type MockKafkaEvent struct{}

func (e MockKafkaEvent) String() string {
	return "A mock event"
}

type MockProducer struct {
	producedCount    int
	outEventsChannel chan kafka.Event
}

func (p *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	p.producedCount += 1
	p.outEventsChannel <- MockKafkaEvent{}

	return nil
}

func (p MockProducer) Events() chan kafka.Event {
	return p.outEventsChannel
}

func NewMockProducer() MockProducer {
	return MockProducer{0, make(chan kafka.Event)}
}

func TestProduceIsCalled(t *testing.T) {
	mockProducer := NewMockProducer()
	client := NewKafkaClient("testing-topic", &mockProducer)
	eventsChan := make(chan simulation.Event)
	var wg sync.WaitGroup

	go client.Start(eventsChan)

	wg.Add(2)

	go func() {
		defer wg.Done()
		eventsChan <- simulation.CreateRandomEvent("ABCD", []string{"sign_up"})
	}()

	go func() {
		defer wg.Done()

		ev := <-mockProducer.outEventsChannel
		fmt.Println(ev)
		assert.Equal(t, mockProducer.producedCount, 1, "Produced count is not one")
	}()

	wg.Wait()
}
