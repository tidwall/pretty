package pretty

import (
	"bytes"
	gojson "encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func j(json interface{}) string {
	var v interface{}
	if err := gojson.Unmarshal([]byte(fmt.Sprintf("%s", json)), &v); err != nil {
		panic(err)
	}
	data, err := gojson.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

var example1 = []byte(`
{
	"name": {
		"last": "Sanders",
		"first": "Janet"
	}, 
	"children": [
		"Andy", "Carol", "Mike"
	],
	"values": [
		10.10, true, false, null, "hello", {}
	],
	"values2": {},
	"values3": [],
	"deep": {"deep":{"deep":[1,2,3,4,5]}}
}
`)

var example2 = `[ 0, 10, 10.10, true, false, null, "hello \" "]`

func TestPretty(t *testing.T) {
	pretty := Pretty(Ugly(Pretty([]byte(example1))))
	assert.Equal(t, j(pretty), j(pretty))
	assert.Equal(t, j(example1), j(pretty))
	pretty = Pretty(Ugly(Pretty([]byte(example2))))
	assert.Equal(t, j(pretty), j(pretty))
	assert.Equal(t, j(example2), j(pretty))
	pretty = Pretty([]byte(" "))
	assert.Equal(t, "", string(pretty))
	pretty = Pretty([]byte("{  "))
	assert.Equal(t, "{", string(pretty))
}

func TestUgly(t *testing.T) {
	ugly := Ugly([]byte(example1))
	var buf bytes.Buffer
	err := gojson.Compact(&buf, []byte(example1))
	assert.Equal(t, nil, err)
	assert.Equal(t, buf.Bytes(), ugly)
	ugly = UglyInPlace(ugly)
	assert.Equal(t, nil, err)
	assert.Equal(t, buf.Bytes(), ugly)
}

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++ {
		b := make([]byte, 1024)
		rand.Read(b)
		Pretty(b)
		Ugly(b)
	}
}

func BenchmarkPretty(t *testing.B) {
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Pretty(example1)
	}
}

func BenchmarkUgly(t *testing.B) {
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Ugly(example1)
	}
}

func BenchmarkUglyInPlace(t *testing.B) {
	example2 := []byte(string(example1))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		UglyInPlace(example2)
	}
}
func BenchmarkJSONIndent(t *testing.B) {
	var dst bytes.Buffer
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		gojson.Indent(&dst, example1, "", "  ")
	}
}

func BenchmarkJSONCompact(t *testing.B) {
	var dst bytes.Buffer
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		gojson.Compact(&dst, example1)
	}
}
