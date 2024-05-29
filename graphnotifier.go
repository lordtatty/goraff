package goraff

import (
	"sync"
)

type StateChangeNotification struct {
	NodeID string
}

type GraphNotifier struct {
	mu        sync.Mutex
	callbacks []func(StateChangeNotification)
}

// Register adds a new callback function.
func (n *GraphNotifier) Listen(callback func(StateChangeNotification)) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.callbacks == nil {
		n.callbacks = make([]func(StateChangeNotification), 0)
	}
	if callback == nil {
		return
	}
	n.callbacks = append(n.callbacks, callback)
}

// Notify triggers all registered callbacks with the given notification.
func (n *GraphNotifier) Notify(notification StateChangeNotification) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, callback := range n.callbacks {
		callback(notification)
	}
}
