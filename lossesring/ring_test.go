package lossesring_test

import (
	"log"
	"minecraft_searcher/lossesring"
	"testing"
)

func TestRing(t *testing.T) {
	r := lossesring.New[int](10)
	for i := 0; i < 12; i++ {
		r.Push(i)
	}
	if r.GetArray()[0] != 2 {
		log.Fatalf("First value in array is %d, but should be 2", r.GetArray()[0])
	}
}
