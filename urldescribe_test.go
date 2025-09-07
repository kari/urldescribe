package urldescribe

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestEmptyURL(t *testing.T) {
	resp, _ := DescribeURL(ctx, "")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestLocalhostURL(t *testing.T) {
	resp, _ := DescribeURL(ctx, "http://localhost/secrets.php")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestRelativeURL(t *testing.T) {
	resp, _ := DescribeURL(ctx, "/url")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}

func TestSchemelessURL(t *testing.T) {
	resp, _ := DescribeURL(ctx, "www.google.com/url")
	if resp != "" {
		t.Errorf("Expected empty string, got %s", resp)
	}
}
