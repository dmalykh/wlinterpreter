package slice

import (
	"github.com/dmalykh/wlinterpreter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	var s = NewStack[int32](10)

	assert.Equal(t, 0, s.GetPosition())
	assert.ErrorIs(t, s.Move(44), wlinterpreter.ErrStackOverflowed)
	assert.ErrorIs(t, s.Move(-44), wlinterpreter.ErrStackOverflowed)
	assert.NoError(t, s.Move(4))
	assert.NoError(t, s.SetValue(8))
	assert.NoError(t, s.Move(1))
	assert.Equal(t, 5, s.GetPosition())
	{
		val, err := s.GetValue(s.GetPosition())
		assert.NoError(t, err)
		assert.Equal(t, int32(0), val)
	}
	assert.NoError(t, s.Move(-1))
	{
		val, err := s.GetValue(s.GetPosition())
		assert.NoError(t, err)
		assert.Equal(t, int32(8), val)
	}
	{
		_, err := s.GetValue(44)
		assert.ErrorIs(t, err, wlinterpreter.ErrStackOverflowed)
	}
	assert.NoError(t, s.SetPosition(3))
	assert.ErrorIs(t, s.SetPosition(323), wlinterpreter.ErrStackOverflowed)

}
