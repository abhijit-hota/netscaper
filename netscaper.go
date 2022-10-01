package netscaper

import (
	"errors"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	href         = "HREF"
	icon         = "ICON"
	iconURI      = "ICON_URI"
	tags         = "TAGS"
	addDate      = "ADD_DATE"
	lastModified = "LAST_MODIFIED"
	private      = "PRIVATE"
	lastVisited  = "LAST_VISITED"
)

var (
	h3                   = regexp.MustCompile("<H3.*>(.*)</H3>")
	h3End                = regexp.MustCompile(`</DL>\s*<p>\s*(<HR>)?$`)
	spaces               = regexp.MustCompile(`\s{2,}`)
	anchorAttributeRegex = regexp.MustCompile(`(HREF|ADD_DATE|LAST_MODIFIED|ICON_URI|ICON|TAGS|PRIVATE|LAST_VISITED)*="(.*?)"`)
	anchorTitle          = regexp.MustCompile(`<A.*>(.*)</A>`)
)

type Bookmark struct {
	Href         string
	Title        string
	Description  string
	Icon         string
	IconURI      string
	Tags         []string
	FolderPath   string
	AddDate      time.Time
	LastModified time.Time
	LastVisited  time.Time
	Private      bool
}

type Options struct {
	// The string used for separating folders in the folder path.
	// Default is ␝
	FolderPathSeparator string

	// Whether to parse description or not.
	// This refers to the <DD> tags
	ParseDescription bool

	// Whether to ignore bookmarklets or not.
	// This refers to URLs starting with javascript:
	IgnoreBookmarklets bool
}

func Parse(str string, opts *Options) ([]Bookmark, error) {
	if opts == nil {
		opts = &Options{
			FolderPathSeparator: "␝",
			ParseDescription:    false,
			IgnoreBookmarklets:  true,
		}
	}

	if !strings.HasPrefix(str, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		return nil, errors.New("not a valid file")
	}
	folderPath := make([]string, 0)

	entities := strings.Split(str, "<DT>")
	bookmarks := make([]Bookmark, 0)

	for _, entity := range entities {
		entity = cleanStr(entity)

		// An H3 tag means a start of a folder
		if strings.HasPrefix(strings.ToUpper(entity), "<H3") {
			res := strings.TrimSpace(h3.FindStringSubmatch(entity)[1])
			folderPath = append(folderPath, res)
		}

		// An single bookmark is a link
		if strings.HasPrefix(strings.ToUpper(entity), "<A") {
			aTag := entity

			if !strings.HasSuffix(entity, "</A>") {
				lastA := strings.LastIndex(entity, "</A>")
				aTag = entity[:lastA+4]
			}
			bm := getAnchorAttributes(aTag)

			// Ignore or keep bookmarklets according to options
			isBookmarklet := strings.HasPrefix(bm.Href, "javascript:")
			if isBookmarklet && opts.IgnoreBookmarklets {
				continue
			}

			// If there's a description it starts with a <DD> tag
			if opts.ParseDescription {
				checkDesc := strings.Split(entity, "<DD>")
				if len(checkDesc) == 2 {
					bm.Description = cleanStr(checkDesc[1])
				}
			}

			bm.FolderPath = strings.Join(folderPath, opts.FolderPathSeparator)
			bookmarks = append(bookmarks, bm)
		}

		// A closing </DL><p> tag means the end of a folder
		if h3End.MatchString(entity) && len(folderPath) > 0 {
			folderPath = folderPath[:len(folderPath)-1]
		}
	}
	return bookmarks, nil
}

func ParseFromFile(path string, opts *Options) ([]Bookmark, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}

	return Parse(string(contents), opts)
}

func getAnchorAttributes(anchorStr string) Bookmark {
	bm := Bookmark{}
	bm.Title = html.UnescapeString(anchorTitle.FindStringSubmatch(anchorStr)[1])
	bm.Tags = make([]string, 0)

	attributeKeyValues := anchorAttributeRegex.FindAllStringSubmatch(anchorStr, -1)

	for _, v := range attributeKeyValues {
		key, value := v[1], v[2]
		switch key {
		case href:
			bm.Href = value
		case icon:
			bm.Icon = value
		case iconURI:
			bm.IconURI = value
		case tags:
			bm.Tags = strings.Split(value, ",")
		case addDate:
			intVal, _ := strconv.Atoi(value)
			bm.AddDate = time.Unix(int64(intVal), 0)
		case lastModified:
			intVal, _ := strconv.Atoi(value)
			bm.LastModified = time.Unix(int64(intVal), 0)
		case lastVisited:
			intVal, _ := strconv.Atoi(value)
			bm.LastVisited = time.Unix(int64(intVal), 0)
		case private:
			bm.Private = value == "1"
		}
	}

	return bm
}

func cleanStr(str string) string {
	return strings.TrimSpace(
		spaces.ReplaceAllString(
			strings.ReplaceAll(str, "\n", ""),
			" ",
		),
	)
}
