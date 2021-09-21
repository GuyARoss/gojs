package ipc

import (
	"testing"
)

func TestIPC(t *testing.T) {
	ipc, err := New("test123")
	if err != nil {
		t.Error("error occurred when creating ipc", err.Error())
		return
	}

	opName := "testOP"
	data := "plzwork"

	done := make(chan bool)

	go func() {
		l, err := ipc.Listen()
		if err != nil {
			t.Error("ipc listener threw error")
		}

		s := <-l.IncomingData

		if s.Data != data && s.OpName != opName {
			t.Errorf("got data=%s, opname=%s \nexpected data=%s, opname=%s",
				s.Data, s.OpName, data, opName,
			)
		}

		done <- true
	}()

	go func() {
		ipc.Write(opName, data)
	}()

	<-done
}
