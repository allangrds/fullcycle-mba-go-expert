package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

type TestEventHandler struct{}

func (h *TestEventHandler) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()

	suite.handler = TestEventHandler{}
	suite.handler2 = TestEventHandler{}
	suite.handler3 = TestEventHandler{}

	suite.event = TestEvent{
		Name:    "TestEvent",
		Payload: "Payload",
	}
	suite.event2 = TestEvent{
		Name:    "TestEvent2",
		Payload: "Payload2",
	}
}

func (suite *EventDispatcherTestSuite) TestEvenDispatcher_Register() {
	assert.True(suite.T(), true)
}

// SetupSuite é executado uma vez antes de todos os testes
func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
