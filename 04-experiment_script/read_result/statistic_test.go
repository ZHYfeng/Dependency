package main

import (
	"fmt"
	"testing"
)

func TestCheckPath(t *testing.T) {
	t.Parallel()
	test := &prevalent{}

	t.Run(fmt.Sprint(), func(t *testing.T) {
		test.average([]*prevalent{})
	})

}
