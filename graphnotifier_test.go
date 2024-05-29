package goraff_test

import (
	"sync"
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestGraphNotifier_ListenAndNotify(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Define a callback function and a variable to capture the notification
	var receivedNotification goraff.StateChangeNotification
	callback := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		receivedNotification = notification
	}

	// Register the callback
	wg.Add(1)
	notifier.Listen(callback)

	// Create a notification
	expectedNotification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks
	notifier.Notify(expectedNotification)

	// Wait for the callback to be invoked
	wg.Wait()

	// Assert that the callback was invoked with the expected notification
	assert.Equal(t, expectedNotification, receivedNotification)
}

func TestGraphNotifier_MultipleCallbacks(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Define callback functions and variables to capture notifications
	receivedNotifications := make([]goraff.StateChangeNotification, 2)
	callback1 := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		receivedNotifications[0] = notification
	}
	callback2 := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		receivedNotifications[1] = notification
	}

	// Register the callbacks
	wg.Add(2)
	notifier.Listen(callback1)
	notifier.Listen(callback2)

	// Create a notification
	expectedNotification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks
	notifier.Notify(expectedNotification)

	// Wait for the callbacks to be invoked
	wg.Wait()

	// Assert that both callbacks were invoked with the expected notification
	assert.Equal(t, expectedNotification, receivedNotifications[0])
	assert.Equal(t, expectedNotification, receivedNotifications[1])
}

func TestGraphNotifier_NoCallbacks(t *testing.T) {
	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Create a notification
	notification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify without any registered callbacks (should not panic)
	assert.NotPanics(t, func() {
		notifier.Notify(notification)
	})
}

func TestGraphNotifier_NilCallback(t *testing.T) {
	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Register a nil callback, this should not be added to the list of callbacks
	notifier.Listen(nil)

	// Create a notification
	notification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks (should not panic even if nil has been passed as a callback)
	assert.NotPanics(t, func() {
		notifier.Notify(notification)
	})
}

func TestGraphNotifier_Concurrency(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	const numCallbacks = 100
	receivedNotifications := make([]goraff.StateChangeNotification, numCallbacks)
	callback := func(i int) func(goraff.StateChangeNotification) {
		return func(notification goraff.StateChangeNotification) {
			defer wg.Done()
			receivedNotifications[i] = notification
		}
	}

	// Register multiple callbacks
	for i := 0; i < numCallbacks; i++ {
		wg.Add(1)
		notifier.Listen(callback(i))
	}

	// Create a notification
	expectedNotification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks
	notifier.Notify(expectedNotification)

	// Wait for all callbacks to be invoked
	wg.Wait()

	// Assert that all callbacks were invoked with the expected notification
	for i := 0; i < numCallbacks; i++ {
		assert.Equal(t, expectedNotification, receivedNotifications[i])
	}
}

func TestGraphNotifier_EmptyNotification(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Define a callback function and a variable to capture the notification
	var receivedNotification goraff.StateChangeNotification
	callback := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		receivedNotification = notification
	}

	// Register the callback
	wg.Add(1)
	notifier.Listen(callback)

	// Create an empty notification
	expectedNotification := goraff.StateChangeNotification{}

	// Notify all registered callbacks
	notifier.Notify(expectedNotification)

	// Wait for the callback to be invoked
	wg.Wait()

	// Assert that the callback was invoked with the expected notification
	assert.Equal(t, expectedNotification, receivedNotification)
}

func TestGraphNotifier_SameCallbackMultipleTimes(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Define a callback function and a counter to track how many times it's called
	callCount := 0
	callback := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		callCount++
	}

	// Register the same callback multiple times
	const timesToRegister = 3
	for i := 0; i < timesToRegister; i++ {
		wg.Add(1)
		notifier.Listen(callback)
	}

	// Create a notification
	notification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks
	notifier.Notify(notification)

	// Wait for all callbacks to be invoked
	wg.Wait()

	// Assert that the callback was called the expected number of times
	assert.Equal(t, timesToRegister, callCount)
}

func TestGraphNotifier_CallbackOrder(t *testing.T) {
	var wg sync.WaitGroup

	// Create a new GraphNotifier instance
	notifier := &goraff.GraphNotifier{}

	// Define callback functions and an array to capture the order of calls
	callOrder := []int{}
	var mu sync.Mutex
	callback1 := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		mu.Lock()
		callOrder = append(callOrder, 1)
		mu.Unlock()
	}
	callback2 := func(notification goraff.StateChangeNotification) {
		defer wg.Done()
		mu.Lock()
		callOrder = append(callOrder, 2)
		mu.Unlock()
	}

	// Register the callbacks
	wg.Add(2)
	notifier.Listen(callback1)
	notifier.Listen(callback2)

	// Create a notification
	notification := goraff.StateChangeNotification{NodeID: "node123"}

	// Notify all registered callbacks
	notifier.Notify(notification)

	// Wait for all callbacks to be invoked
	wg.Wait()

	// Assert that the callbacks were called in the order they were registered
	expectedOrder := []int{1, 2}
	assert.Equal(t, expectedOrder, callOrder)
}
