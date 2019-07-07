// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import discovery "bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
import generation "bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
import logrus "github.com/sirupsen/logrus"
import mock "github.com/stretchr/testify/mock"
import model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// GenerateSpecificationTestCases provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Generator) GenerateSpecificationTestCases(_a0 *logrus.Entry, _a1 generation.GeneratorConfig, _a2 discovery.ModelDiscovery, _a3 *model.Context) generation.SpecRun {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 generation.SpecRun
	if rf, ok := ret.Get(0).(func(*logrus.Entry, generation.GeneratorConfig, discovery.ModelDiscovery, *model.Context) generation.SpecRun); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(generation.SpecRun)
	}

	return r0
}
