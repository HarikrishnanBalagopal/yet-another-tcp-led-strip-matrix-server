package main

import (
	"net"
	"time"

	"github.com/crazy3lf/colorconv"
	"github.com/sirupsen/logrus"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func fill_leds(leds []byte, t float32) {
	const num_rows = 13
	const num_cols = num_rows
	const mid_row = num_rows / 2
	const mid_col = num_cols / 2
	// var cen_x = 0.5*float32(math.Cos(float64(t))) + 0.5
	// const cen_y = 0.0
	// const cir_r = 0.5
	// const cir_r2 = cir_r * cir_r
	for row := 0; row < num_rows; row++ {
		for _col := 0; _col < num_cols; _col++ {
			col := _col
			if row%2 == 1 {
				col = num_cols - 1 - _col
			}
			idx := 3 * (row*num_cols + _col)
			row_d := row - mid_row
			if row_d < 0 {
				row_d = -row_d
			}
			col_d := col - mid_col
			if col_d < 0 {
				col_d = -col_d
			}
			rect_id := row_d
			if col_d > row_d {
				rect_id = col_d
			}

			// x := float32(col) / float32(num_cols)
			// y := float32(row) / float32(num_rows)
			// dx := (x - cen_x)
			// dy := (y - cen_y)
			// dist2 := dx*dx + dy*dy
			red, green, blue, err := colorconv.HSLToRGB(float64((int(t*25)+rect_id*30)%360), 1.0, 0.5)
			must(err)
			// var red byte = 0
			// if dist2 <= cir_r2 {
			// 	red = 255
			// }
			leds[idx+0] = red
			leds[idx+1] = green
			leds[idx+2] = blue
		}
	}
}

func main() {
	logrus.Infof("main start")
	remoteAddr := net.TCPAddr{
		// IP:   net.IPv4(192, 168, 0, 8),
		IP:   net.IPv4(192, 168, 0, 5),
		Port: 8888,
	}
	conn, err := net.DialTCP("tcp", nil, &remoteAddr)
	must(err)
	defer conn.Close()

	const n_total = 3 * 13 * 13
	leds := make([]byte, n_total)
	var t float32 = 0.0

	for {
		fill_leds(leds, t)
		t += 0.1
		// logrus.Infof("leds: %+v", leds)
		{
			_, err := conn.Write([]byte("frame1234\n"))
			must(err)
		}
		n, err := conn.Write(leds)
		must(err)
		if n != n_total {
			logrus.Fatalf("expected: %d actual n: %d", n_total, n)
		}
		time.Sleep(25 * time.Millisecond)
	}
	// logrus.Infof("main end")
}
