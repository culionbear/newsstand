package newsstand

import "unsafe"

type Slot[T any] struct {
	ch chan T
}

func (s *Signal[T]) NewSlot(size uint) *Slot[T] {
	if size == 0 {
		size = s.size
	}
	slot := &Slot[T]{
		ch: make(chan T, size),
	}
	s.slots.Store(uintptr(unsafe.Pointer(slot)), slot)
	return slot
}

func (s *Slot[T]) Subscribe() <-chan T {
	return s.ch
}

func (s *Slot[T]) close() {
	close(s.ch)
}
