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
	ipc.Dispose()
}
func TestIPC_UnsyncProcessCall(t *testing.T) {
	ipc, err := New("test123")
	if err != nil {
		t.Error("error occurred when creating ipc", err.Error())
		return
	}

	opName := "testOP"
	data := "plzwork"

	done := make(chan bool)

	go func() {
		ipc.Write(opName, data)
	}()

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

	<-done
	ipc.Dispose()
}

func TestIPC_SeparateInstancesSamePipe(t *testing.T) {
	ipc1, err := New("test123")
	if err != nil {
		t.Error("error occurred when creating ipc 1", err.Error())
		return
	}
	ipc2, err := New("test123")
	if err != nil {
		t.Error("error occurred when creating ipc 2", err.Error())
		return
	}

	opName := "testOP"
	data := "plzwork"

	done := make(chan bool)

	go func() {
		l, err := ipc2.Listen()
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
		ipc1.Write(opName, data)
	}()

	<-done
	ipc1.Dispose()
	ipc2.Dispose()
}

func TestNormalizeDataOut(t *testing.T) {
	nOut := normalizeDataOut(`
	<html>
    <manifest>

    </manifest>
    <style>
        .center {
            text-align: center;
        }
    </style>
    <body>        
        <div class="center">
            <h1 class="ui"> Settings </h1>  
        </div>
    </body>
</html>
`)
	if nOut != ` <html> <manifest> </manifest> <style> .center { text-align: center; } </style> <body> <div class="center"> <h1 class="ui"> Settings </h1> </div> </body> </html> ` {
		t.Errorf("normalization response not expected")
	}
}
