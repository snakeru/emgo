package main

// Gauge is ready to used gauge that can be updated using
type Gauge struct {
	n    int16
	half int16
	min  int16
	max  int16
}

func (g *Gauge) SetMin(min int) {
	g.min = int16(min)
}

func (g *Gauge) SetMax(max int) {
	g.max = int16(max)
}

// HandleHalf treats m as fixed-point number with 1-bit fractional part. It
// adjust m by previoussly discarded fraction and returns integer part of m.
func (g *Gauge) handleHalf(m int) int {
	if m&1 != 0 {
		if g.half != 0 {
			m += int(g.half)
			g.half = 0
		} else {
			if m > 0 {
				g.half = 1
			} else {
				g.half = -1
			}
			m -= int(g.half)
		}
	}
	return m / 2
}

func (g *Gauge) Val() int {
	return int(g.n)
}

func (g *Gauge) SetVal(n int) {
	switch {
	case n > int(g.max):
		n = int(g.max)
	case n < int(g.min):
		n = int(g.min)
	}
	g.n = int16(n)
}

func (g *Gauge) Add(m int) {
	m = g.handleHalf(m)
	g.SetVal(int(g.n) + m)
}

func (g *Gauge) AddCube(m int) {
	m = g.handleHalf(m)
	g.SetVal(int(g.n) + m*m*m)
}
