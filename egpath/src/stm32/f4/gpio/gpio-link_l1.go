// Package gpio allows setup and use GPIO ports.
package gpio

import "unsafe"

// Common for STM32F4xx and STM32L1xx

// GPIO represents registers of one GPIO port
type Port struct {
	mode   uint32 `C:"volatile"`
	otype  uint32 `C:"volatile"`
	ospeed uint32 `C:"volatile"`
	pupd   uint32 `C:"volatile"`
	id     uint32 `C:"volatile"`
	od     uint32 `C:"volatile"`
	bsr    uint32 `C:"volatile"`
	lck    uint32 `C:"volatile"`
	afl    uint32 `C:"volatile"`
	afh    uint32 `C:"volatile"`
}

const (
	base uintptr = 0x40020000
	step uintptr = 0x400
)

var (
	A = (*Port)(unsafe.Pointer(base + step*0))
	B = (*Port)(unsafe.Pointer(base + step*1))
	C = (*Port)(unsafe.Pointer(base + step*2))
	D = (*Port)(unsafe.Pointer(base + step*3))
	E = (*Port)(unsafe.Pointer(base + step*4))
	F = (*Port)(unsafe.Pointer(base + step*5))
	G = (*Port)(unsafe.Pointer(base + step*6))
	H = (*Port)(unsafe.Pointer(base + step*7))
	I = (*Port)(unsafe.Pointer(base + step*8))
	J = (*Port)(unsafe.Pointer(base + step*9))
	K = (*Port)(unsafe.Pointer(base + step*10))
)

// Number returns port number.
func (g *Port) Number() int {
	return int((uintptr(unsafe.Pointer(g)) - base) / step)
}

type Mode byte

const (
	In Mode = iota
	Out
	Alt
	Analog
)

// Mode returns I/O mode for n-th bit
func (g *Port) Mode(n uint) Mode {
	n *= 2
	return Mode(g.mode>>n) & 3
}

// SetMode sets I/O mode for n-th bit
func (g *Port) SetMode(n uint, mode Mode) {
	n *= 2
	g.mode = g.mode&^(3<<n) | uint32(mode)<<n
}

type OutType byte

const (
	PushPull OutType = iota
	OpenDrain
)

// OutType returns current type of n-th output bit
func (g *Port) OutType(n uint) OutType {
	return OutType(g.otype>>n) & 1
}

// SetOuttype sets type for n-th output bit
func (g *Port) SetOutType(n uint, ot OutType) {
	g.otype = g.otype&^(1<<n) | uint32(ot)<<n
}

type Speed byte

// OutSpeed return current speed for n-th output bit
func (g *Port) OutSpeed(n uint) Speed {
	n *= 2
	return Speed(g.ospeed>>n) & 3
}

// SetOutSpeed sets speed for n-th output bit
func (g *Port) SetOutSpeed(n uint, speed Speed) {
	n *= 2
	g.ospeed = g.ospeed&^(3<<n) | uint32(speed)<<n
}

type Pull byte

const (
	NoPull = iota
	PullUp
	PullDown
)

// Pull returns current pull state of of n-th output bit.
func (g *Port) Pull(n uint) Pull {
	n *= 2
	return Pull(g.pupd>>n) & 3
}

// SetPull sets internal pull-up/pull-down cirquits for n-th output bit.
func (g *Port) SetPull(n uint, pull Pull) {
	n *= 2
	g.pupd = g.pupd&^(3<<n) | uint32(pull)<<n
}

// SetBit sets n-th output bit to 1
func (g *Port) SetBit(n uint) {
	g.bsr = uint32(1) << n
}

// ClearBit sets n-th output bit to 0
func (g *Port) ClearBit(n uint) {
	g.bsr = uint32(0x10000) << n
}

// SetBits sets output bits on positions specified by bits to 1
func (g *Port) SetBits(bits uint16) {
	g.bsr = uint32(bits)
}

// ClearBits sets output bits on positions specified by bits to 0
func (g *Port) ClearBits(bits uint16) {
	g.bsr = uint32(bits) << 16
}

// ClearAndSet sets whole BSRR register.
// High 16 bits in bsr specifies which bits should be 0.
// Low 16 bits in bss specifies which bits should be 1.
// Set bits has priority above clear bits.
func (g *Port) ClearAndSet(bsr uint32) {
	g.bsr = bsr
}

func (g *Port) Store(bits uint16) {
	g.od = uint32(bits)
}

func (g *Port) Load() uint16 {
	return uint16(g.id)
}

func (g *Port) Bit(n uint) bool {
	return g.id&(uint32(1)<<n) != 0
}

// AltFunc returns current alternate function for n-th bit in port g.
func (g *Port) AltFunc(n uint) AltFunc {
	var af uint32
	if n < 8 {
		af = g.afl
	} else {
		af = g.afh
		n -= 8
	}
	n *= 4
	return AltFunc(af>>n) & 0xf
}

// SetAltFunc sets alternate function af for n-th bit in port g.
func (g *Port) SetAltFunc(n uint, af AltFunc) {
	n *= 4
	if n < 32 {
		g.afl = g.afl&^(0xf<<n) | uint32(af)<<n
	} else {
		n -= 32
		g.afh = g.afh&^(0xf<<n) | uint32(af)<<n
	}
}