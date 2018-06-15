package router

import (
	"testing"
)

func TestRouter_getRoot(t *testing.T) {
	testRouter := &Router{}
	testZone := &Zone{
		parent: &testRouter.Zone,
		root:   testRouter,
	}

	if testZone.Root() != testRouter {
		t.Errorf("The test zone parent must be the test router: %p %p", testZone.Root(), testRouter)
	}
}
