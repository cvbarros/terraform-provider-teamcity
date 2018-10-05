package teamcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//PropertyAssertions is a helper for tests that are defined in teamcity package. Not sure how to factor these out in a proper place.
type PropertyAssertions struct {
	a *assert.Assertions
	t *testing.T
}

func newPropertyAssertions(t *testing.T) *PropertyAssertions {
	return &PropertyAssertions{a: assert.New(t), t: t}
}

func (p *PropertyAssertions) assertPropertyValue(props *Properties, name string, value string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	if v, ok := propMap[name]; ok {
		p.a.Equal(value, v)
	} else {
		p.a.Contains(propMap, name)
	}
}

func (p *PropertyAssertions) assertPropertyDoesNotExist(props *Properties, name string) {
	require.NotNil(p.t, props)

	propMap := props.Map()

	p.a.NotContains(propMap, name)
}
