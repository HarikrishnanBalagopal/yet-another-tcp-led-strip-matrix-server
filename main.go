package main

import (
	"errors"
	"net"
	"net/http"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const expected_length = 13 * 13 * 4
const rgb_frame_length = 13 * 13 * 3

var rgb_frame = make([]byte, rgb_frame_length)

func setup() {
	for i := 0; i < rgb_frame_length; i++ {
		if i%2 == 0 {
			rgb_frame[i] = 0
		} else {
			rgb_frame[i] = 255
		}
	}
}

func main() {
	logrus.Infof("main start")
	setup()

	logrus.Infof("tcp setup start")
	remoteAddr := net.TCPAddr{
		IP: net.IPv4(192, 168, 0, 8),
		// IP:   net.IPv4(192, 168, 0, 5),
		Port: 8888,
	}
	conn, err := net.DialTCP("tcp", nil, &remoteAddr)
	must(err)
	defer conn.Close()

	frame_chan := make(chan interface{}, 10)
	frame_chan <- nil

	go func() {
		logrus.Infof("waiting for frames in another thread")
		for {
			<-frame_chan // wait for a frame to arrive
			// logrus.Infof("got a new frame: %+v", rgb_frame)
			// logrus.Infof("leds: %+v", leds)
			_, err := conn.Write([]byte("frame1234\n"))
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					logrus.Errorf("failed to write frame1234 it's a  write: broken pipe error, retry connection")
					conn.Close() // Is this required?
					conn, err = net.DialTCP("tcp", nil, &remoteAddr)
					if err != nil {
						logrus.Fatalf("failed to reconnect over TCP. error: %q", err)
					}
					// defer conn.Close()
					_, err := conn.Write([]byte("frame1234\n"))
					must(err)
				} else {
					logrus.Fatalf("failed to write to the TCP connection. error: %q", err)
				}
			}
			n, err := conn.Write(rgb_frame)
			must(err)
			if n != rgb_frame_length {
				logrus.Fatalf("expected: %d actual n: %d", rgb_frame_length, n)
			}
			logrus.Infof("sent the new frame over the TCP connection")
			// time.Sleep((1000 / 60) * time.Millisecond)
		}
	}()
	logrus.Infof("tcp setup end")

	router := mux.NewRouter()
	if router == nil {
		panic("nil router")
	}
	router.PathPrefix("/ws").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Errorf("failed to upgrade to a websocket")
			return
		}
		conn.WriteJSON(map[string]string{"msg": "hello world from the web socket server"})
		conn.WriteMessage(websocket.BinaryMessage, []byte("this is a binary message from the server"))
		for {
			msgtype, frame, err := conn.ReadMessage()
			if err != nil {
				logrus.Errorf("failed to read a message from the websocket. error: %q", err)
				return
			}
			if msgtype != websocket.BinaryMessage {
				logrus.Errorf("expected a binary message on the websocket. actual: '%s'", string(frame))
				continue
			}
			if len(frame) != expected_length {
				logrus.Errorf("expected length '%d'. actual length: '%d'", expected_length, len(frame))
				continue
			}
			// logrus.Infof("got a binary message on the websocket of length: %d\n%+v", len(frame), frame)
			for row := 0; row < 13; row++ {
				for col := 0; col < 13; col++ {
					src := 4 * (row*13 + col)
					col_alt := col
					if row%2 == 1 {
						col_alt = 13 - 1 - col
					}
					dst := 3 * (row*13 + col_alt)
					rgb_frame[dst+0] = frame[src+0]
					rgb_frame[dst+1] = frame[src+1]
					rgb_frame[dst+2] = frame[src+2]
				}
			}
			// for src, dst := 0, 0; dst < rgb_frame_length; src, dst = src+4, dst+3 {
			// 	rgb_frame[dst+0] = frame[src+0]
			// 	rgb_frame[dst+1] = frame[src+1]
			// 	rgb_frame[dst+2] = frame[src+2]
			// }
			frame_chan <- nil
		}
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	logrus.Infof("listening on port 8080 -> http://127.0.0.1:8080/")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
	logrus.Infof("main end")
}
