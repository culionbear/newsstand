package newsstand_test

import (
	"github.com/culionbear/newsstand"
	"testing"
	"time"
)

func s1(signal *newsstand.Signal[bool], t *testing.T) {
	slot := signal.NewSlot(0)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-slot.Subscribe():
			return
		case <-ticker.C:
			t.Log(1)
		}
	}
}

func s2(signal *newsstand.Signal[bool], t *testing.T) {
	slot := signal.NewSlot(0)
	ticker := time.NewTicker(time.Second)
	n := 0
	for {
		select {
		case ok := <-slot.Subscribe():
			if ok {
				return
			}
		case <-ticker.C:
			t.Log(2)
			n++
			if n == 1 {
				signal.CloseSlot(slot)
			}
		}
	}
}

func TestSignal(t *testing.T) {
	signal := newsstand.New[bool](0)
	signal.Connect(func(bool) {
		t.Log("3")
	})
	go s1(signal, t)
	go s2(signal, t)
	time.Sleep(time.Second * 5)
	signal.Publish(true)
	//signal.Close()
	time.Sleep(time.Second * 5)
}
