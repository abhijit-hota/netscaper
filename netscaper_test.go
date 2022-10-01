package netscaper

import (
	"q"
	"testing"
)

func TestBasicParser(t *testing.T) {
	bookmarks, err := Parse("<!DOCTYPE NETSCAPE-Bookmark-file-1>ffsdffsdfsdf", nil)
	if err != nil {
		panic(err)
	}
	// testing.
	q.Q(bookmarks) // DEBUG
}
