package dtrain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	EmptyHolder = LayoutHolder{}
)

func TestBaseLayoutLoad(t *testing.T) {

	container := Load("tests/full")

	assert.NotEqual(t, EmptyHolder, container.Get("index"))
	assert.NotEqual(t, EmptyHolder, container.Get("list"))
	assert.NotEqual(t, EmptyHolder, container.Get("single"))
	assert.NotEqual(t, EmptyHolder, container.Get("form"))

	assert.Equal(t, EmptyHolder, container.Get("test"))
	assert.Equal(t, EmptyHolder, container.Get("head"))
	assert.Equal(t, EmptyHolder, container.Get("foot"))
	assert.Equal(t, EmptyHolder, container.Get("baseof"))
}

func TestContentSpecifiedLayoutLoad(t *testing.T) {

	container := Load("tests/full")

	assert.NotEqual(t, EmptyHolder, container.Get("product/index"))
	assert.NotEqual(t, EmptyHolder, container.Get("product/list"))
	assert.NotEqual(t, EmptyHolder, container.Get("product/single"))
	assert.NotEqual(t, EmptyHolder, container.Get("product/form"))

	assert.NotEqual(t, EmptyHolder, container.Get("user/index"))
	assert.NotEqual(t, EmptyHolder, container.Get("user/list"))
	assert.NotEqual(t, EmptyHolder, container.Get("user/single"))
	assert.NotEqual(t, EmptyHolder, container.Get("user/form"))
}
