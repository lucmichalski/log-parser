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
	assert.Equal(t, stat.ElapsedTime, 0)
	elapsed := stat.Stop()
	assert.NotEqual(t, stat.ElapsedTime, 0)
	assert.Equal(t, stat.ElapsedTime, elapsed)
}

func TestAddLog(t *testing.T) {
	stat := NewStat("request", "")
	assert.Empty(t, stat.Logs)
	file := "/var/log/nginx/access.log"
	stat.AddLog(file)
	assert.Equal(t, stat.Logs, []string{file})
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

func TestEmptyRegexp(t *testing.T) {
	stat := NewStat("request", "")
	assert.Nil(t, stat.GroupByRegexp)
}

func TestRegexp(t *testing.T) {
	exp := `^\w+\s+(\S+)(?:\?|$)`
	stat := NewStat("request", exp)
	assert.Equal(t, stat.GroupByRegexp.String(), exp)
}

func TestGroupByRegexp(t *testing.T) {
	stat := NewStat("request", `^\w+\s+(\S+)`)
	uri := "/foo/bar"
	request := "GET " + uri
	entry := &gonx.Entry{"request": request}

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)

	counter, ok := stat.Data[request]
	assert.False(t, ok)
	assert.Equal(t, counter, 0)

	// Uri should be used as data key because we have regexp to extract it
	counter, ok = stat.Data[uri]
	assert.True(t, ok)
	assert.Equal(t, counter, 1)
}

func TestBadRegexp(t *testing.T) {
	// Invalid Regexp required request to be numeric field
	stat := NewStat("request", `^(\d+)$`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
}

func TestNoSubmatchRegexp(t *testing.T) {
	// Invalid Regexp required request to be numeric field
	stat := NewStat("request", `^\w+`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)

	counter, ok := stat.Data[request]
	assert.False(t, ok)
	assert.Equal(t, counter, 0)

	// Request method was used for grouping
	counter, ok = stat.Data["GET"]
	assert.True(t, ok)
	assert.Equal(t, counter, 1)
}
