package main_test

import (
	"testing"

	lookout "github.com/jakealves/atlas-game-lookout"
	"github.com/stretchr/testify/assert"
)

func TestParseMapTileCommand(t *testing.T) {
	label, zoom, long, lat, shift := lookout.ParseMapTileCommand("!maptile long:-60.12 lat:53.12")
	assert.Equal(t, "", label, "they should be equal")
	assert.Equal(t, 64, zoom, "they should be equal")
	assert.Equal(t, -60.12, long, "they should be equal")
	assert.Equal(t, 53.12, lat, "they should be equal")
	assert.Equal(t, 0, shift, "they should be equal")

	label, zoom, long, lat, shift = lookout.ParseMapTileCommand("!maptile zoom:54 shift:43 label:\"fart\"")
	assert.Equal(t, "fart", label, "they should be equal")
	assert.Equal(t, 54, zoom, "they should be equal")
	assert.Equal(t, 0.0, long, "they should be equal")
	assert.Equal(t, 0.0, lat, "they should be equal")
	assert.Equal(t, 43, shift, "they should be equal")
}
