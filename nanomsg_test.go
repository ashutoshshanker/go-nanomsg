// Go binding for nanomsg

package nanomsg

import (
	"bytes"
	"testing"
	"time"
)


func TestReqReqp(t *testing.T) {
	var err error
	var rep, req *Socket
	socketAddress := "inproc://a"

	if rep, err = NewSocket(SP, REP); err != nil {
		t.Fatal(err)
	}
	if _, err = rep.Bind(socketAddress); err != nil {
		t.Fatal(err)
	}
	if req, err = NewSocket(SP, REQ); err != nil {
		t.Fatal(err)
	}
	if _, err = req.Connect(socketAddress); err != nil {
		t.Fatal(err)
	}

	if err = req.Send([]byte("ABC"), 0); err != nil {
		t.Fatal(err)
	}
	if data, err := rep.Recv(0); err != nil {
		t.Fatal(err)
	} else if (bytes.Compare(data, []byte("ABC")) != 0) {
		t.Errorf("Unexpected data received: %s", data)
	}

	if err = rep.Close(); err != nil {
		t.Fatal(err)
	}
	if err = req.Close(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkInprocThroughput(b *testing.B) {
	b.StopTimer()

	worker := func() {
		var err error
		var s *Socket
		if s, err = NewSocket(SP, PAIR); err != nil {
			b.Fatal(err)
		}
		if _, err = s.Connect("inproc://inproc_bench"); err != nil {
			b.Fatal(err)
		}

		for i := 0; i < b.N; i++ {
			if data, err := s.Recv(0); err != nil {
				b.Fatal(err)
			} else if err = s.Send(data, 0); err != nil {
				b.Fatal(err)
			}
		}

		if err = s.Close(); err != nil {
			b.Fatal(err)
		}
	}

	var err error
	var s *Socket
	if s, err = NewSocket(SP, PAIR); err != nil {
		b.Fatal(err)
	}
	if _, err = s.Bind("inproc://inproc_bench"); err != nil {
		b.Fatal(err)
	}

	// Wait a bit till the worker routine blocks in Recv().
	// TODO signal the worker to die somehow
	var s2 *Socket
	fixme := false
	if fixme {
		go worker()
		time.Sleep(0 * 100 * time.Nanosecond)
	} else {
		if s2, err = NewSocket(SP, PAIR); err != nil {
			b.Fatal(err)
		}
		if _, err = s2.Connect("inproc://inproc_bench"); err != nil {
			b.Fatal(err)
		}
	}

	buf := bytes.Repeat([]byte{111}, 10240)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err = s.Send(buf, 0); err != nil {
			b.Fatal(err)
		}
		if !fixme {
			if data, err := s2.Recv(0); err != nil {
				b.Fatal(err)
			} else if err = s2.Send(data, 0); err != nil {
				b.Fatal(err)
			}
		}
		if _, err := s.Recv(0); err != nil {
			b.Fatal(err)
		}
	}

	if !fixme {
		if err = s2.Close(); err != nil {
			b.Fatal(err)
		}
	}
	if err = s.Close(); err != nil {
		b.Fatal(err)
	}
}

func ExampleBind(t *testing.T) {
	socket, err := NewSocket(SP, REP)
	if err != nil {
		panic(err)
	}
	socket.Bind("inproc://a")
}
