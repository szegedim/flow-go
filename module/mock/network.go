// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"

	network "github.com/onflow/flow-go/network"
)

// Network is an autogenerated mock type for the Network type
type Network struct {
	mock.Mock
}

// Register provides a mock function with given fields: channel, engine
func (_m *Network) Register(channel network.Channel, engine network.Engine) (network.Conduit, error) {
	ret := _m.Called(channel, engine)

	var r0 network.Conduit
	if rf, ok := ret.Get(0).(func(network.Channel, network.Engine) network.Conduit); ok {
		r0 = rf(channel, engine)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conduit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(network.Channel, network.Engine) error); ok {
		r1 = rf(channel, engine)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
