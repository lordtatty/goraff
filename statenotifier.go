package goraff

import (
	"sync"
)

type StateChangeNotification struct {
	NodeID string
}

type StateNotifier struct {
	mu        sync.Mutex
	callbacks []func(StateChangeNotification)
}

// Register adds a new callback function.
func (n *StateNotifier) Listen(callback func(StateChangeNotification)) {
	if n.callbacks == nil {
		n.callbacks = make([]func(StateChangeNotification), 0)
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.callbacks = append(n.callbacks, callback)
}

// Notify triggers all registered callbacks with the given notification.
func (n *StateNotifier) Notify(notification StateChangeNotification) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, callback := range n.callbacks {
		callback(notification)
	}
}
