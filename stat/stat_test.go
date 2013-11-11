package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstructor(t *testing.T) {
	stat := NewStat("request", "")
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.GroupBy, "request")
}

func TestTiming(t *testing.T) {
	start := time.Now()
	stat := NewStat("request", "")
	assert.WithinDuration(t, start, stat.StartedAt, time.Duration(time.Millisecond),
		"Constructor should setup StartedAt")
	assert.Empty(t, stat.Logs)
	assert.Equal(t, stat.ElapsedTime, 0)
	stat.Stop()
	assert.NotEqual(t, stat.ElapsedTime, 0)
}

func TestEmptyRegexp(t *testing.T) {
	stat := NewStat("request", "")
	assert.Nil(t, stat.GroupByRegexp)
}

func TestRegexp(t *testing.T) {
	exp := `^\w+\s+(\S+)(?:\?|$)`
	stat := NewStat("request", exp)
	assert.Equal(t, stat.GroupByRegexp.String(), exp)
}

func TestAdd(t *testing.T) {
	stat := NewStat("request", "")
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)
	counter, ok := stat.Data[request]
	assert.True(t, ok)
	assert.Equal(t, counter, 1)

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 2)
	counter, ok = stat.Data[request]
	assert.True(t, ok)
	assert.Equal(t, counter, 2)
}

func TestAddInvalid(t *testing.T) {
	stat := NewStat("request", "")
	entry := &gonx.Entry{"foo": "bar"}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, len(stat.Data), 0)
}

func TestAddLog(t *testing.T) {
	stat := NewStat("request", "")
	file := "/var/log/nginx/access.log"
	stat.AddLog(file)
	assert.Equal(t, stat.Logs, []string{file})
}
