package main

import (
	"errors"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"

	"q"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(path string) string {
	arr, err := os.ReadFile(path)
	must(err)

	return string(arr)
}

var h3 *regexp.Regexp = regexp.MustCompile("<H3.*>(.*)</H3>")
var h3End *regexp.Regexp = regexp.MustCompile(`</DL>\s*<p>\s*(<HR>)?$`)
var spaces = regexp.MustCompile(`\s{2,}`)
var anchorAttributeRegex = regexp.MustCompile(`(HREF|ADD_DATE|LAST_MODIFIED|ICON_URI|ICON|TAGS)*="(.*?)"`)
var anchorTitle = regexp.MustCompile(`<A.*>(.*)</A>`)

type Bookmark struct {
	Href         string
	Title        string
	Icon         string
	IconUri      string
	Tags         []string
	FolderPath   string
	AddDate      int64
	LastModified int64
}

const (
	HREF          = "HREF"
	ICON          = "ICON"
	ICON_URI      = "ICON_URI"
	TAGS          = "TAGS"
	ADD_DATE      = "ADD_DATE"
	LAST_MODIFIED = "LAST_MODIFIED"
)

func ParseAnchor(anchorStr string) Bookmark {
	bm := Bookmark{}
	bm.Title = html.UnescapeString(anchorTitle.FindStringSubmatch(anchorStr)[1])
	bm.Tags = make([]string, 0)

	attributeKeyValues := anchorAttributeRegex.FindAllStringSubmatch(anchorStr, -1)

	for _, v := range attributeKeyValues {
		key, value := v[1], v[2]
		switch key {
		case HREF:
			bm.Href = value
		case ICON:
			bm.Icon = value
		case ICON_URI:
			bm.IconUri = value
		case TAGS:
			bm.Tags = strings.Split(value, ",")
		case ADD_DATE:
			intVal, _ := strconv.Atoi(value)
			bm.AddDate = int64(intVal)
		case LAST_MODIFIED:
			intVal, _ := strconv.Atoi(value)
			bm.LastModified = int64(intVal)
		}
	}

	return bm
}

func Parse(str string) ([]Bookmark, error) {

	if !strings.HasPrefix(str, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		return nil, errors.New("nope")
	}
	path := make([]string, 0)

	entities := strings.Split(str, "<DT>")
	bookmarks := make([]Bookmark, 0)

	for _, entity := range entities {
		entity = strings.TrimSpace(spaces.ReplaceAllString(strings.ReplaceAll(entity, "\n", ""), " "))

		if strings.HasPrefix(strings.ToUpper(entity), "<H3") {
			res := strings.TrimSpace(h3.FindStringSubmatch(entity)[1])
			path = append(path, res)
		}

		if strings.HasPrefix(strings.ToUpper(entity), "<A") {
			aTag := entity
			if !strings.HasSuffix(entity, "</A>") {
				lastA := strings.LastIndex(entity, "</A>")
				aTag = entity[:lastA+4]
			}
			bm := ParseAnchor(aTag)
			bm.FolderPath = strings.Join(path, "/")
			bookmarks = append(bookmarks, bm)
		}

		if h3End.MatchString(entity) && len(path) > 0 {
			path = path[:len(path)-1]
		}
	}
	return bookmarks, nil
}

func main() {
	str := ReadFile("latest-bm.html")
	bookmarks, _ := Parse(str)
	q.Q(len(bookmarks)) // DEBUG
}
