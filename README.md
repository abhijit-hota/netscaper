[![Go Reference](https://pkg.go.dev/badge/github.com/abhijit-hota/netscaper.svg)](https://pkg.go.dev/github.com/abhijit-hota/netscaper)

## Netscaper

## What is it?

A Netscape Bookmark File parser. In case you don't know what that is, it's the `bookmarks.html` browsers blurt out when you **`Export bookmarks as HTML`**.

**This library takes such files and outputs an array of bookmarks.**

### Features 

- ✅ Parse bookmarked items
  - ✅ Title, link, tags, date added, last modified
  - ✅ Descriptions inside `<DD>` 
  - ✅ `PRIVATE` attributes
  - ✅ `LAST_VISITED` attributes
- ✅ Parse Folders  
  - ✅ outputs the full folder path when in flattened mode
  - 🚧 exposes functions to output a folder tree  
  - 🚧 handles folder with the same name and level

### Usage

```go
package main

import (
	"github.com/abhijit-hota/netscaper"
)

func main() {
    // If you have a file
    pathToFile := "./bookmarks.html"
	books, err := netscaper.ParseFromFile(pathToFile, nil /* default options */)
	if err != nil {
		panic(err)
	}

    // If you have a string with the file contents:
    contents := "..." 
  	books, err := netscaper.Parse(contents, nil /* default options */)
	if err != nil {
		panic(err)
	} 
}
```
Check the [API reference](https://pkg.go.dev/github.com/abhijit-hota/netscaper) for more.

### Notes & References

The code uses regexes to parse HTML. Sigh. I might rewrite this using the [`html`](https://pkg.go.dev/golang.org/x/net/html) package or using [lexical scanning](https://youtu.be/HxaD_trXwRE) but this was written for a really specific purpose and here we are.

- https://github.com/kafene/netscape-bookmark-parser
- https://learn.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/aa753582(v=vs.85)?redirectedfrom=MSDN



> *15 years ago this library would be all the rage.*
> 
> \- Me