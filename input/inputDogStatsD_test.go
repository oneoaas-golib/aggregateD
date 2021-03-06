package input

import (
	"testing"
)

func TestInvalidMessage(t *testing.T) {
	message := "foo:5|gauge324324"
	result, err := parseDogStatsDMetric(message)

	if err == nil {
		t.Error("expected error, got", result)
	}

	message = "foogauge324324"
	result, err = parseDogStatsDMetric(message)

	if err == nil {
		t.Error("expected error, got", result)
	}
}

func TestMetricParse(t *testing.T) {
	//message := "metric.name:value|type|@sample_rate|#tag1:value,tag2"
	message := "foo:5|gauge|@0.5|#nonkvtag,tag1:firstvalue,tag2:second,tag3:third"
	result, _ := parseDogStatsDMetric(message)

	if result.Name != "foo" {
		t.Error("Expected name foo got", result.Name)
	}

	if result.Value != 5 {
		t.Error("Expected value of 5 got", result.Value)
	}

	if result.Sampling != 0.5 {
		t.Error("Exected sampling of 0.5 got", result.Sampling)
	}

	v1 := result.Tags["tag2"]
	if v1 != "second" {
		t.Error("value of tag1 was expected to be second got", v1)
	}

	v3 := result.Tags["tag3"]
	if v3 != "third" {
		t.Error("value of tag3 was expected to be third got", v3)
	}

	v2 := result.Tags["nonkvtag"]
	if v2 != "nonkvtag" {
		t.Error("value of nonkvtag was expected to be nonkvtag got", v2)
	}
}

// func TestEventParse(t *testing.T) {
// 	message := "_e{title.length,text.length}:title|text|d:date_happened|h:hostname|p:priority|t:alert_type|#tag1,tag2"
// 	parseDogStatsDEvent(message)
//
// }
