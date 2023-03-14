package expect

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// Expectation is a embedding type for assert.Assertions
type Expectation struct {
	*assert.Assertions
}

// New makes a new Expectation object for the specified TestingT.
func New(t assert.TestingT) *Expectation {
	assert := assert.New(t)
	return &Expectation{assert}
}

// TODO(james): Equality depends on order of repeated fields. This may break some tests unnecessarily.
// ProtoEqual asserts that the specified protobuf messages are equal.
func (a *Expectation) ProtoEqual(expected, actual proto.Message) bool {
	return a.True(
		proto.Equal(expected, actual),
		fmt.Sprintf("These two protobuf messages are not equal:\nexpected: %v\nactual:  %v", expected, actual),
	)
}
