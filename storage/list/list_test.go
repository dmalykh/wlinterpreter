package list

import (
	"github.com/dmalykh/wlinterpreter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListStorage_Alloc(t *testing.T) {
	var ls = ListStorage{
		data: make(map[string]*element),
	}

	assert.NoError(t, ls.Alloc(`chumbawamba`))
	val, exists := ls.data[`chumbawamba`]
	assert.True(t, exists)
	assert.Nil(t, val)
}

func TestListStorage_Empty(t *testing.T) {
	var ls = New()

	assert.True(t, ls.Empty(`hocus`))
	assert.NoError(t, ls.Alloc(`pocus`))
	assert.True(t, ls.Empty(`pocus`))
	assert.NoError(t, ls.Append(`pocus`, 42))
	assert.False(t, ls.Empty(`pocus`))
}

func TestListStorage_Append(t *testing.T) {
	var ls = ListStorage{
		data: make(map[string]*element),
	}

	assert.ErrorIs(t, ls.Append(`doors`, 31), wlinterpreter.ErrNotAllocatedBlock)

	assert.NoError(t, ls.Alloc(`doors`))
	assert.NoError(t, ls.Append(`doors`, 42))
	assert.Equal(t, 42, ls.data[`doors`].val)
	assert.NoError(t, ls.Append(`doors`, 49))
	assert.Equal(t, 49, ls.data[`doors`].val)
	assert.Equal(t, 42, ls.data[`doors`].prev.val)

}

func TestListStorage_Pop(t *testing.T) {
	var ls = ListStorage{
		data: make(map[string]*element),
	}

	// Check not allocated error
	{
		val, err := ls.Pop(`nazareth`)
		assert.ErrorIs(t, err, wlinterpreter.ErrNotAllocatedBlock)
		assert.Equal(t, 0, val)
	}

	// Check add data
	assert.NoError(t, ls.Alloc(`nazareth`))
	assert.NoError(t, ls.Append(`nazareth`, 99))
	assert.NoError(t, ls.Append(`nazareth`, 77))

	// Check pop
	{
		val, err := ls.Pop(`nazareth`)
		assert.NoError(t, err)
		assert.Equal(t, 77, val)
	}
	{
		val, err := ls.Pop(`nazareth`)
		assert.NoError(t, err)
		assert.Equal(t, 99, val)
	}
	{
		val, err := ls.Pop(`nazareth`)
		assert.ErrorIs(t, err, wlinterpreter.ErrNoValue)
		assert.Equal(t, 0, val)
	}
	// And check impossible remove error
	{
		assert.ErrorIs(t, ls.remove(`nazareth`, nil), ErrNothingToRemove)
		assert.ErrorIs(t, ls.remove(`fgkdyweuvi`, &element{}), wlinterpreter.ErrNotAllocatedBlock)
	}
}

func TestListStorage_Last(t *testing.T) {
	var ls = New()
	{
		val, err := ls.Last(`fgkdyweuvi`)
		assert.ErrorIs(t, err, wlinterpreter.ErrNotAllocatedBlock)
		assert.Equal(t, 0, val)
	}
	assert.NoError(t, ls.Alloc(`ramjam`))
	{
		val, err := ls.Last(`ramjam`)
		assert.ErrorIs(t, err, wlinterpreter.ErrNoValue)
		assert.Equal(t, 0, val)
	}
	assert.NoError(t, ls.Append(`ramjam`, 500))
	{
		val, err := ls.Last(`ramjam`)
		assert.NoError(t, err)
		assert.Equal(t, 500, val)
	}
}

func TestListStorage_Update(t *testing.T) {
	var ls = New()
	assert.ErrorIs(t, ls.Update(`kiss`, 32), wlinterpreter.ErrNotAllocatedBlock)
}
