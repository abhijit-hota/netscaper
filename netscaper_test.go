package netscaper

import (
	"testing"
)

func TestBasicParser(t *testing.T) {
	_, err := Parse("<!DOCTYPE NETSCAPE-Bookmark-file-1>ffsdffsdfsdf", nil)
	if err != nil {
		panic(err)
	}
}
