package core

import (
	"math/rand/v2"
	"sync"
	"testing"
)

func TestRand(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 1000; i++ {
			RandSrc.With(func(v **rand.Rand) {
				(*v).Int64()
			})
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			RandSrc.With(func(v **rand.Rand) {
				(*v).Int64()
			})
		}
		wg.Done()
	}()
	wg.Wait()
}
