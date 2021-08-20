package pulsar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBBPulsarJob(t *testing.T) {
	j := NewBBPulsarJob("myBBPulsarJob", "myAppId")
	assert.Equal(t, "custom", j.TypeID, "Wrong typeid")
	assert.Equal(t, "myBBPulsarJob", j.Name, "Wrong name")
	assert.Equal(t, "myAppId", j.AppID, "Wrong appid")
}

func TestNewJSPulsarJob(t *testing.T) {
	j := NewJSPulsarJob("myJSPulsarJob", "myAppId", "myHost", "myURLPath")
	assert.Equal(t, "latency", j.TypeID, "Wrong typeid")
	assert.Equal(t, "myJSPulsarJob", j.Name, "Wrong name")
	assert.Equal(t, "myAppId", j.AppID, "Wrong appid")
	assert.Equal(t, "myHost", j.Config.Host, "Wrong host")
	assert.Equal(t, "myURLPath", j.Config.URL_Path, "Wrong url_path")
}
