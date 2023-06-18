package router

import (
	"testing"
)

func TestRouter_getRoot(t *testing.T) {
	testRouter := &routerImp{}
	testZone := &zone{
		parent: &testRouter.zone,
		root:   testRouter,
	}

	if testZone.Root() != testRouter {
		t.Errorf("The test zone parent must be the test router: %p %p", testZone.Root(), testRouter)
	}
}
