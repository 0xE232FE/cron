package core

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/alaingilbert/cron/internal/mtx"
	mathRand "math/rand/v2"
)

var RandSrc mtx.Mtx[*mathRand.Rand]

func init() {
	var seedBytes [16]byte
	_, _ = rand.Read(seedBytes[:])
	seedState := binary.LittleEndian.Uint64(seedBytes[:8])
	seedStream := binary.LittleEndian.Uint64(seedBytes[8:])
	RandSrc.Set(mathRand.New(mathRand.NewPCG(seedState, seedStream)))
}
