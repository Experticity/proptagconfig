package proptagconfig_test

import (
	"testing"

	. "github.com/Experticity/proptagconfig"
	"github.com/Experticity/tagconfig"
	"github.com/magiconair/properties"
	"github.com/stretchr/testify/assert"
)

type HackerReference struct {
	Handle        string `props:"handle"`
	ZeroCool      bool   `props:"zero-cool"`
	ComputerCount int    `props:"computer-count"`
}

func TestPropTagConfigGet(t *testing.T) {
	p := properties.NewProperties()

	handle := "the.plague"
	_, _, err := p.Set("handle", handle)
	assert.NoError(t, err)

	pg := &PropTagConfig{p}

	hr := &HackerReference{}
	err = tagconfig.Process(pg, hr)
	assert.NoError(t, err)

	assert.Equal(t, handle, hr.Handle)
}

func TestPropTagConfigGetMissing(t *testing.T) {
	p := properties.NewProperties()

	pg := &PropTagConfig{p}

	hr := &HackerReference{}
	err := tagconfig.Process(pg, hr)
	assert.NoError(t, err)

	assert.Empty(t, hr.Handle)
}

func TestPropSet(t *testing.T) {
	p := properties.NewProperties()

	pt := &PropTagConfig{p}

	hr := &HackerReference{
		Handle:        "the.plague",
		ZeroCool:      true,
		ComputerCount: 17,
	}
	err := tagconfig.PopulateExternalSource(pt, hr)
	assert.NoError(t, err)
	assert.Equal(t, hr.Handle, p.MustGet("handle"))
	assert.True(t, p.GetBool("zero-cool", false))
	assert.Equal(t, hr.ComputerCount, p.GetInt("computer-count", 1))
}

func TestPropSetNonString(t *testing.T) {
	type AnotherHackerReference struct {
		Murphy struct{ Learn string } `props:"zero-cool"`
	}

	pt := &PropTagConfig{properties.NewProperties()}

	hr := &AnotherHackerReference{struct{ Learn string }{"Revenge"}}
	err := tagconfig.PopulateExternalSource(pt, hr)
	assert.EqualError(t, err, ErrUnsupportedSetType.Error())
}

func TestPopulatePropertiesFromStruct(t *testing.T) {
	tests := []struct {
		setup     func() interface{}
		expFunc   func() *properties.Properties
		shouldErr bool
	}{
		{
			setup: func() interface{} {
				return &struct {
					Name string `props:"name"`
				}{"Lebowski"}
			},
			expFunc: func() *properties.Properties {
				p := properties.NewProperties()
				_, _, err := p.Set("name", "Lebowski")
				assert.NoError(t, err)
				return p
			},
		},

		{
			setup: func() interface{} {
				return nil
			},

			shouldErr: true,
		},
	}

	for _, tt := range tests {
		p, err := PopulatePropertiesFromStruct(tt.setup())
		if tt.shouldErr {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, tt.expFunc(), p)
	}
}
