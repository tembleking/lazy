package lazy

import (
	"sync"
)

type Lazy[T any] struct {
	value *T
	mutex sync.RWMutex
}

func (l *Lazy[T]) GetOrInit(init func() T) T {
	l.mutex.RLock()
	if l.value != nil {
		defer l.mutex.RUnlock()
		return *l.value
	}
	l.mutex.RUnlock()

	l.mutex.Lock()
	defer l.mutex.Unlock()
	// we duplicate the check in here because only the first mutex will
	// succeed adding value to l.value, the rest of goroutines will just
	// return the value assigned
	if l.value != nil {
		return *l.value
	}

	initValue := init()
	l.value = &initValue
	return *l.value
}

func (l *Lazy[T]) GetOrTryInit(init func() (T, error)) (T, error) {
	l.mutex.RLock()
	if l.value != nil {
		defer l.mutex.RUnlock()
		return *l.value, nil
	}
	l.mutex.RUnlock()

	l.mutex.Lock()
	defer l.mutex.Unlock()
	// we duplicate the check in here because only the first mutex will
	// succeed adding value to l.value, the rest of goroutines will just
	// return the value assigned
	if l.value != nil {
		return *l.value, nil
	}

	initValue, err := init()
	if err != nil {
		return initValue, err
	}

	l.value = &initValue
	return *l.value, nil
}
