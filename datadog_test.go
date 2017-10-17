package datadog

import "testing"

func Test_parseMetricName(t *testing.T) {
	expectName := "com.example.metricName"
	expectTags := []string{
		"tag1:value1",
		"tag2:va-lue2",
		"ta-g",
		"job:hello/world.sql",
	}
	ddMetric := parseMetricName("com.example.metricName[tag1:value1,tag2:va-lue2,ta-g,job:hello/world.sql]")

	if ddMetric.name != expectName {
		t.Logf("Expected metric name %v but got %v", expectName, ddMetric.name)
		t.Fail()
	}

	tags := ddMetric.tags
	if len(tags) != len(expectTags) {
		t.Logf("Expected to have %v tags parsed but found %v", len(expectTags), len(tags))
		t.Fail()
	}

	if len(tags) == 3 {
		for i := range tags {
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

func Test_baseTags(t *testing.T) {

	b := baseTags(Config{})

	if len(b) != 0 {
		t.Log("Expected length of base tags to be zero")
		t.Fail()
	}

	b2 := baseTags(Config{Environment: "development"})
	if len(b2) != 1 {
		t.Log("Expected to have exactly one base tag")
		t.Fail()
	}

	if b2[0] != "environment:development" {
		t.Log("Expected tag to be 'environment:development'")
		t.Fail()
	}

	b3 := baseTags(Config{
		AppName:     "metrics",
		Environment: "development",
	})

	if len(b3) != 2 {
		t.Log("Expected to have exactly two tags")
		t.Fail()
	}

	if b3[1] != "app:metrics" {
		t.Log("Expected tag to be 'app:metrics'")
		t.Fail()
	}
}
