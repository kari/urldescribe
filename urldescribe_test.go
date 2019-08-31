package urldescribe

import "testing"

func TestEmptyURL(t *testing.T) {
	resp := DescribeURL("")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestLocalhostURL(t *testing.T) {
	resp := DescribeURL("http://localhost/secrets.php")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestRelativeURL(t *testing.T) {
	resp := DescribeURL("/url")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestSchemelessURL(t *testing.T) {
	resp := DescribeURL("www.google.com/url")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}
