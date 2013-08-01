package opc

// Raver plaid
//   A rainbowy pattern with moving diagonal black stripes

import (
	"github.com/longears/pixelslinger/colorutils"
	"github.com/longears/pixelslinger/midi"
	"math"
	"time"
)

func MakePatternRaverPlaid(locations []float64) ByteThread {
	return func(bytesIn chan []byte, bytesOut chan []byte, midiState *midi.MidiState) {

		var (
			// how many sine wave cycles are squeezed into our n_pixels
			// 24 happens to create nice diagonal stripes on the wall layout
			freq_r float64 = 24
			freq_g float64 = 24
			freq_b float64 = 24

			// how many seconds the color sine waves take to shift through a complete cycle
			speed_r float64 = 7
			speed_g float64 = -13
			speed_b float64 = 19
		)

		for bytes := range bytesIn {
			n_pixels := len(bytes) / 3
			t := float64(time.Now().UnixNano())/1.0e9 - 9.4e8

			// fill in bytes array
			for ii := 0; ii < n_pixels; ii++ {
				//--------------------------------------------------------------------------------

				pct := float64(ii) / float64(n_pixels)

				// replicate a quirk in the original python version of this pattern
				pct /= 2

				// diagonal black stripes
				pct_jittered := colorutils.PosMod2((pct * 77), 37)
				blackstripes := colorutils.Cos(pct_jittered, t*0.05, 1, -1.5, 1.5) // offset, period, minn, maxx
				blackstripes_offset := colorutils.Cos(t, 0.9, 60, -0.5, 3)
				blackstripes = colorutils.Clamp(blackstripes+blackstripes_offset, 0, 1)

				// 3 sine waves for r, g, b which are out of sync with each other
				r := blackstripes * colorutils.Remap(math.Cos((t/speed_r+pct*freq_r)*math.Pi*2), -1, 1, 0, 1)
				g := blackstripes * colorutils.Remap(math.Cos((t/speed_g+pct*freq_g)*math.Pi*2), -1, 1, 0, 1)
				b := blackstripes * colorutils.Remap(math.Cos((t/speed_b+pct*freq_b)*math.Pi*2), -1, 1, 0, 1)

				bytes[ii*3+0] = colorutils.FloatToByte(r)
				bytes[ii*3+1] = colorutils.FloatToByte(g)
				bytes[ii*3+2] = colorutils.FloatToByte(b)

				//--------------------------------------------------------------------------------
			}
			bytesOut <- bytes
		}
	}
}
