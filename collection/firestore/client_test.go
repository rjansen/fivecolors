package firestore

import (
	"testing"

	"github.com/rjansen/l"
	firestore "github.com/rjansen/raizel/firestore/mock"
)

type clientTest struct {
	logger       l.Logger
	client       *firestore.ClientMock
	setupTest    func(t *testing.T, test *clientTest)
	tearDownTest func(t *testing.T, test *clientTest)
}

func (test *clientTest) setup(t *testing.T) {
	test.logger = l.LoggerDefault
	test.client = firestore.NewClientMock()
	if test.setupTest != nil {
		test.setupTest(t, test)
	}
}

func (test *clientTest) tearDown(t *testing.T) {
	if test.tearDownTest != nil {
		test.tearDownTest(t, test)
	}
	// test.client.Close()
}
