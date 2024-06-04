package notifiers

import (
	"sync"

	"github.com/lordtatty/goraff"
)

type GraphNotifier struct {
	mu        sync.Mutex
	callbacks []func(goraff.GraphChangeNotification)
}

// Register adds a new callback function.
func (n *GraphNotifier) Listen(callback func(goraff.GraphChangeNotification)) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.callbacks == nil {
		n.callbacks = make([]func(goraff.GraphChangeNotification), 0)
	}
	if callback == nil {
		return
	}
	n.callbacks = append(n.callbacks, callback)
}

// Notify triggers all registered callbacks with the given notification.
func (n *GraphNotifier) Notify(notification goraff.GraphChangeNotification) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, callback := range n.callbacks {
		callback(notification)
	}
}
