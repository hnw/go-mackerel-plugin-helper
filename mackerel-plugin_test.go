package mackerelplugin

import (
	"bytes"
	"math"
	"testing"
	"time"
)

func TestCalcDiff(t *testing.T) {
	val1 := 10.0
	val2 := 0.0
	now := time.Now()
	last := time.Unix(now.Unix()-10, 0)

	diff, err := calcDiff(val1, now, val2, last)
	if diff != 60 {
		t.Errorf("calcDiff: %f should be %f", diff, 60.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestCalcDiffWithReset(t *testing.T) {
	val := 10.0
	now := time.Now()
	lastval := 12345.0
	last := time.Unix(now.Unix()-60, 0)

	diff, err := calcDiff(val, now, lastval, last)
	if err == nil {
		t.Errorf("calcDiffUInt32 with counter reset should cause an error: %f", diff)
	}
}

func TestCalcDiffWithUInt32WithReset(t *testing.T) {
	val := uint32(10)
	now := time.Now()
	lastval := uint32(12345)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := calcDiffUInt(uint64(val), now, uint64(lastval), last, 10, math.MaxUint32)
	if err == nil {
		t.Errorf("calcDiffUInt32 with counter reset should cause an error: %f", diff)
	}
}

func TestCalcDiffWithUInt32Overflow(t *testing.T) {
	val := uint32(10)
	now := time.Now()
	lastval := math.MaxUint32 - uint32(10)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := calcDiffUInt(uint64(val), now, uint64(lastval), last, 10, math.MaxUint32)
	if diff != 21.0 {
		t.Errorf("calcDiff: last: %d, now: %d, %f should be %f", val, lastval, diff, 21.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestCalcDiffWithUInt64WithReset(t *testing.T) {
	val := uint64(10)
	now := time.Now()
	lastval := uint64(12345)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := calcDiffUInt(val, now, lastval, last, 10, math.MaxUint64)
	if err != nil {
	} else {
		t.Errorf("calcDiffUInt64 with counter reset should cause an error: %f", diff)
	}
}

func TestCalcDiffWithUInt64Overflow(t *testing.T) {
	val := uint64(10)
	now := time.Now()
	lastval := math.MaxUint64 - uint64(10)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := calcDiffUInt(val, now, lastval, last, 10, math.MaxUint64)
	if diff != 21.0 {
		t.Errorf("calcDiff: last: %d, now: %d, %f should be %f", val, lastval, diff, 21.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestPrintValueUint32(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", uint32(10), now)

	expected := []byte("test\t10\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func TestPrintValueUint64(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", uint64(10), now)

	expected := []byte("test\t10\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func TestPrintValueFloat64(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", float64(10.0), now)

	expected := []byte("test\t10.000000\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func ExampleFormatValues() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(1000)}
	lastStat := map[string]interface{}{"cmd_get": uint64(500), ".last_diff.cmd_get": 300.0}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.cmd_get	500.000000	1437227240
}

func ExampleFormatValuesWithCounterReset() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(10)}
	lastStat := map[string]interface{}{"cmd_get": uint64(500), ".last_diff.cmd_get": 300.0}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatFloatValuesWithCounterReset() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "float"}
	stat := map[string]interface{}{"cmd_get": 10.0}
	lastStat := map[string]interface{}{"cmd_get": 500.0, ".last_diff.cmd_get": 300.0}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatValuesWithOverflow() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(500)}
	lastStat := map[string]interface{}{"cmd_get": uint64(math.MaxUint64 - 100), ".last_diff.cmd_get": float64(100.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.cmd_get	601.000000	1437227240
}

func ExampleFormatValuesWithOverflowAndTooHighDifference() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(500)}
	lastStat := map[string]interface{}{"cmd_get": uint64(math.MaxUint64 - 100), ".last_diff.cmd_get": float64(10.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatValuesWithOverflowAndNoLastDiff() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(500)}
	lastStat := map[string]interface{}{"cmd_get": uint64(math.MaxUint64 - 100)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatValuesWithWildcard() {
	var mp MackerelPlugin
	prefix := "foo.#"
	metric := Metrics{Name: "bar", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"foo.1.bar": uint64(1000), "foo.2.bar": uint64(2000)}
	lastStat := map[string]interface{}{"foo.1.bar": uint64(500), ".last_diff.foo.1.bar": float64(2.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1.bar	500.000000	1437227240
}

func ExampleFormatValuesWithWildcardAndNoDiff() {
	var mp MackerelPlugin
	prefix := "foo.#"
	metric := Metrics{Name: "bar", Label: "Get", Diff: false}
	stat := map[string]interface{}{"foo.1.bar": float64(1000)}
	lastStat := map[string]interface{}{"foo.1.bar": float64(500), ".last_diff.foo.1.bar": float64(2.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1.bar	1000.000000	1437227240
}

func ExampleFormatValuesWithWildcardAstarisk() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "*", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"foo.1": uint64(1000), "foo.2": uint64(2000)}
	lastStat := map[string]interface{}{"foo.1": uint64(500), ".last_diff.foo.1": float64(2.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1	500.000000	1437227240
}

// an example implementation
type MemcachedPlugin struct {
}

var graphdef map[string](Graphs) = map[string](Graphs){
	"memcached.cmd": Graphs{
		Label: "Memcached Command",
		Unit:  "integer",
		Metrics: [](Metrics){
			Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"},
		},
	},
}

func (m MemcachedPlugin) GraphDefinition() map[string](Graphs) {
	return graphdef
}

func (m MemcachedPlugin) FetchMetrics() (map[string]interface{}, error) {
	var stat map[string]interface{}
	return stat, nil
}

func ExampleOutputDefinitions() {
	var mp MemcachedPlugin
	helper := NewMackerelPlugin(mp)
	helper.OutputDefinitions()

	// Output:
	// # mackerel-agent-plugin
	// {"graphs":{"memcached.cmd":{"label":"Memcached Command","unit":"integer","metrics":[{"name":"cmd_get","label":"Get","type":"uint64","stacked":false,"scale":0}]}}}
}

func TestToUint32(t *testing.T) {
	if ret := toUint32(uint32(100)); ret != uint32(100) {
		t.Errorf("toUint32(uint32) returns incorrect value: %v expected to be %v", ret, uint32(100))
	}

	if ret := toUint32(uint64(100)); ret != uint32(100) {
		t.Errorf("toUint32(uint64) returns incorrect value: %v expected to be %v", ret, uint32(100))
	}

	if ret := toUint32(float64(100)); ret != uint32(100) {
		t.Errorf("toUint32(float64) returns incorrect value: %v expected to be %v", ret, uint32(100))
	}

	if ret := toUint32("100"); ret != uint32(100) {
		t.Errorf("toUint32(string) returns incorrect value: %v expected to be %v", ret, uint32(100))
	}
}

func TestToUint64(t *testing.T) {
	if ret := toUint64(uint32(100)); ret != uint64(100) {
		t.Errorf("toUint64(uint32) returns incorrect value: %v expected to be %v", ret, uint64(100))
	}

	if ret := toUint64(uint64(100)); ret != uint64(100) {
		t.Errorf("toUint64(uint64) returns incorrect value: %v expected to be %v", ret, uint64(100))
	}

	if ret := toUint64(float64(100)); ret != uint64(100) {
		t.Errorf("toUint64(float64) returns incorrect value: %v expected to be %v", ret, uint64(100))
	}

	if ret := toUint64("100"); ret != uint64(100) {
		t.Errorf("toUint64(string) returns incorrect value: %v expected to be %v", ret, uint64(100))
	}
}

func TestToFloat64(t *testing.T) {
	if ret := toFloat64(uint32(100)); ret != float64(100) {
		t.Errorf("toFloat64(uint32) returns incorrect value: %v expected to be %v", ret, float64(100))
	}

	if ret := toFloat64(uint64(100)); ret != float64(100) {
		t.Errorf("toFloat64(uint64) returns incorrect value: %v expected to be %v", ret, float64(100))
	}

	if ret := toFloat64(float64(100)); ret != float64(100) {
		t.Errorf("toFloat64(float64) returns incorrect value: %v expected to be %v", ret, float64(100))
	}

	if ret := toFloat64("100"); ret != float64(100) {
		t.Errorf("toFloat64(string) returns incorrect value: %v expected to be %v", ret, float64(100))
	}
}
