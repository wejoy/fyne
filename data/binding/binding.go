//go:generate go run gen.go

package binding

import "sync"

// DataItem is the base interface for all bindable data items.
//
// Since: 2.0.0
type DataItem interface {
	// AddListener attaches a new change listener to this DataItem.
	// Listeners are called each time the data inside this DataItem changes.
	// Additionally the listener will be triggered upon successful connection to get the current value.
	AddListener(DataListener)
	// RemoveListener will detach the specified change listener from the DataItem.
	// Disconnected listener will no longer be triggered when changes occur.
	RemoveListener(DataListener)
}

// DataListener is any object that can register for changes in a bindable DataItem.
// See NewDataListener to define a new listener using just an inline function.
//
// Since: 2.0.0
type DataListener interface {
	DataChanged()
}

// NewDataListener is a helper function that creates a new listener type from a simple callback function.
//
// Since: 2.0.0
func NewDataListener(fn func()) DataListener {
	return &listener{fn}
}

type listener struct {
	callback func()
}

func (l *listener) DataChanged() {
	l.callback()
}

type base struct {
	listeners []DataListener
	lock      sync.RWMutex
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for _, listen := range b.listeners { // TODO maybe fix? workaround for list not being able to de-register
		if listen == l {
			return
		}
	}

	b.listeners = append(b.listeners, l)
	queueItem(l.DataChanged)
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for i, listen := range b.listeners {
		if listen != l {
			continue
		}

		if i == len(b.listeners)-1 {
			b.listeners = b.listeners[:len(b.listeners)-1]
		} else {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
		}
	}
}

func (b *base) trigger() {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for _, listen := range b.listeners {
		queueItem(listen.DataChanged)
	}
}
