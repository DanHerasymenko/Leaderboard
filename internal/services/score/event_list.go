package score

import "sync"

type EventList struct {
	listeners map[chan<- *Score]struct{}
	mu        sync.RWMutex
	scores    chan Score
}

func StartEventList() *EventList {
	el := &EventList{
		listeners: make(map[chan<- *Score]struct{}),
		scores:    make(chan Score, 100),
	}
	go el.run()
	return el

}

func (el *EventList) PostScore(score *Score) {
	el.scores <- *score
}

func (el *EventList) Subscribe(l chan<- *Score) {
	el.mu.Lock()
	defer el.mu.Unlock()

	el.listeners[l] = struct{}{}
}

func (el *EventList) Unsubscribe(l chan<- *Score) {
	el.mu.Lock()
	defer el.mu.Unlock()

	delete(el.listeners, l)
}

func (el *EventList) run() {
	for score := range el.scores {
		el.mu.RLock()
		for l := range el.listeners {
			l <- &score
		}
		el.mu.RUnlock()
	}
}
