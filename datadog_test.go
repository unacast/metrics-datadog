package datadog

import "testing"

func Test_parseMetricName(t *testing.T) {
	expectName := "com.example.metricName"
	expectTags := []string{
		"tag1:value1",
		"tag2:va-lue2",
		"ta-g",
	}
	ddMetric := parseMetricName("com.example.metricName[tag1:value1,tag2:va-lue2,ta-g]")

	if ddMetric.name != expectName {
		t.Logf("Expected metric name %v but got %v", expectName, ddMetric.name)
		t.Fail()
	}

	tags := ddMetric.tags
	if len(tags) != 3 {
		t.Logf("Expected to have %v tags parsed but found %v", len(expectTags), len(tags))
		t.Fail()
	}

	if len(tags) == 3 {
		for i, _ := range tags {
			if tags[i] != expectTags[i] {

				t.Logf("Exepected first metric tag to have be %v but was %v", tags[i], expectTags[i])
				t.Fail()

			}
		}
	}
}

func Test_parseMetricName_empty_tags(t *testing.T) {
	expectName := "com.example.metricName"

	ddMetric := parseMetricName("com.example.metricName[]")

	if ddMetric.name != expectName {
		t.Logf("Expected metric name %v but got %v", expectName, ddMetric.name)
		t.Fail()
	}

	if len(ddMetric.tags) != 0 {
		t.Logf("Expected to have 0 tags parsed but found %v", len(ddMetric.tags))
		t.Fail()
	}

}
