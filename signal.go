package newsstand

import (
	"sync"
	"unsafe"
)

type Signal[T any] struct {
	size  uint
	ch    chan T
	slots sync.Map
	funcs sync.Map
}

func New[T any](size uint) *Signal[T] {
	s := &Signal[T]{
		size: size,
		ch:   make(chan T, size),
	}
	go s.listen()
	return s
}

func (s *Signal[T]) Publish(msg T) {
	s.ch <- msg
}

func (s *Signal[T]) listen() {
	for msg := range s.ch {
		s.slots.Range(func(key, value any) bool {
			slot := value.(*Slot[T])
			slot.ch <- msg
			return true
		})
	}
}

func (s *Signal[T]) Connect(f func(T)) {
	key := uintptr(unsafe.Pointer(&f))
	if _, ok := s.funcs.Load(key); ok {
		return
	}
	slot := s.NewSlot(0)
	s.funcs.Store(key, slot)
	s.slots.Store(uintptr(unsafe.Pointer(slot)), slot)
	go func() {
		for msg := range slot.Subscribe() {
			f(msg)
		}
	}()
}

func (s *Signal[T]) Disconnect(f func(T)) {
	key := uintptr(unsafe.Pointer(&f))
	if value, ok := s.funcs.LoadAndDelete(key); ok {
		s.CloseSlot(value.(*Slot[T]))
	}
}

func (s *Signal[T]) CloseSlot(slot *Slot[T]) {
	s.slots.Delete(uintptr(unsafe.Pointer(slot)))
	slot.close()
}

func (s *Signal[T]) Close() {
	close(s.ch)
	s.funcs.Range(func(key, value any) bool {
		s.funcs.Delete(key)
		s.CloseSlot(value.(*Slot[T]))
		return true
	})
	s.slots.Range(func(key, value any) bool {
		s.CloseSlot(value.(*Slot[T]))
		return true
	})
}
