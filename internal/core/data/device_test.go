package data

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/edgexfoundry/edgex-go/internal/core/data/messaging"
	"github.com/edgexfoundry/edgex-go/internal/pkg/correlation/models"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"

	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
	"github.com/edgexfoundry/edgex-go/pkg/clients/metadata/mocks"
	"github.com/edgexfoundry/edgex-go/pkg/clients/types"
	contract "github.com/edgexfoundry/edgex-go/pkg/models"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

var testEvent contract.Event
var testRoutes *mux.Router

const (
	testDeviceName string = "Test Device"
	testOrigin     int64  = 123456789
	testBsonString string = "57e59a71e4b0ca8e6d6d4cc2"
	testUUIDString string = "ca93c8fa-9919-4ec5-85d3-f81b2b6a7bc1"
)

// Mock implementation of the event publisher for testing purposes
type mockEventPublisher struct{}

func TestCheckMaxLimit(t *testing.T) {
	reset()

	testedLimit := math.MinInt32

	expectedNil := checkMaxLimit(testedLimit)

	if expectedNil != nil {
		t.Errorf("Should not exceed limit")
	}
}

func TestCheckMaxLimitOverLimit(t *testing.T) {
	reset()

	testedLimit := math.MaxInt32

	expectedErr := checkMaxLimit(testedLimit)

	if expectedErr == nil {
		t.Errorf("Exceeded limit and should throw error")
	}
}

func newMockEventPublisher(config messaging.PubSubConfiguration) messaging.EventPublisher {
	return &mockEventPublisher{}
}

func (zep *mockEventPublisher) SendEventMessage(e models.Event) error {
	return nil
}

func TestMain(m *testing.M) {
	testRoutes = LoadRestRoutes()
	LoggingClient = logger.NewMockClient()
	mdc = newMockDeviceClient()
	ep = newMockEventPublisher(messaging.PubSubConfiguration{})
	chEvents = make(chan interface{}, 10)
	os.Exit(m.Run())
}

// Supporting methods
// Reset() re-initializes dependencies for each test
func reset() {
	Configuration = &ConfigurationStruct{}
	testEvent.ID = testBsonString
	testEvent.Device = testDeviceName
	testEvent.Origin = testOrigin
	testEvent.Readings = buildReadings()
}

func newMockDeviceClient() *mocks.DeviceClient {
	client := &mocks.DeviceClient{}

	mockAddressable := contract.Addressable{
		Address:  "localhost",
		Name:     "Test Addressable",
		Port:     3000,
		Protocol: "http"}

	mockDeviceResultFn := func(id string, ctx context.Context) contract.Device {
		if bson.IsObjectIdHex(id) {
			return contract.Device{Id: id, Name: testDeviceName, Addressable: mockAddressable}
		}
		return contract.Device{}
	}
	client.On("Device", "valid", context.Background()).Return(mockDeviceResultFn, nil)
	client.On("Device", "404", context.Background()).Return(mockDeviceResultFn,
		types.NewErrServiceClient(http.StatusNotFound, []byte{}))
	client.On("Device", mock.Anything, context.Background()).Return(mockDeviceResultFn, fmt.Errorf("some error"))

	mockDeviceForNameResultFn := func(name string, ctx context.Context) contract.Device {
		device := contract.Device{Id: uuid.New().String(), Name: name, Addressable: mockAddressable}

		return device
	}
	client.On("DeviceForName", testDeviceName, context.Background()).Return(mockDeviceForNameResultFn, nil)
	client.On("DeviceForName", "404", context.Background()).Return(mockDeviceForNameResultFn,
		types.NewErrServiceClient(http.StatusNotFound, []byte{}))
	client.On("DeviceForName", mock.Anything, context.Background()).Return(mockDeviceForNameResultFn,
		fmt.Errorf("some error"))

	return client
}

func buildReadings() []contract.Reading {
	ticks := db.MakeTimestamp()
	r1 := contract.Reading{Id: bson.NewObjectId().Hex(),
		Name:     "Temperature",
		Value:    "45",
		Origin:   testOrigin,
		Created:  ticks,
		Modified: ticks,
		Pushed:   ticks,
		Device:   testDeviceName}

	r2 := contract.Reading{Id: bson.NewObjectId().Hex(),
		Name:     "Pressure",
		Value:    "1.01325",
		Origin:   testOrigin,
		Created:  ticks,
		Modified: ticks,
		Pushed:   ticks,
		Device:   testDeviceName}
	readings := []contract.Reading{}
	readings = append(readings, r1, r2)
	return readings
}

func handleDomainEvents(bitEvents []bool, wait *sync.WaitGroup, t *testing.T) {
	until := time.Now().Add(250 * time.Millisecond) //Kill this loop after quarter second.
	for time.Now().Before(until) {
		select {
		case evt := <-chEvents:
			switch evt.(type) {
			case DeviceLastReported:
				e := evt.(DeviceLastReported)
				if e.DeviceName != testDeviceName {
					t.Errorf("DeviceLastReported name mismatch %s", e.DeviceName)
					return
				}
				bitEvents[0] = true
				break
			case DeviceServiceLastReported:
				e := evt.(DeviceServiceLastReported)
				if e.DeviceName != testDeviceName {
					t.Errorf("DeviceLastReported name mismatch %s", e.DeviceName)
					return
				}
				bitEvents[1] = true
				break
			}
		default:
			//	Without a default case in here, the select block will hang.
		}
	}
	wait.Done()
}
