package main

import (
	"LukasDerBaum/templates"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

func md_to_HTML(md string) string {
	var buf bytes.Buffer

    err := goldmark.Convert([]byte(md), &buf)
    if err != nil {
        panic(err)
    }

    var s string = buf.String()
    //fmt.Println(s)
    s = strings.ReplaceAll(s, "&lt;", "<")
    s = strings.ReplaceAll(s, "&gt;", ">")
    // fmt.Println(s)
    
    
    return s
}

func write(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func main() {
	md, err := os.ReadFile("content/index.md")
	if err != nil {
		panic(err)
	}

	htmlBody := md_to_HTML(string(md))
	os.Remove("build/")
	os.Mkdir("build", 0755)
	f, err := os.Create("build/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// fmt.Println("Rendering page...", "My Site", htmlBody)
	// page_body := templ.Raw(htmlBody)
	//fmt.Println(page_body)
	page := templates.Page("My Site", htmlBody)
	fmt.Println(page)
	err = page.Render(context.Background(), f)
	if err != nil {
		panic(err)
	}

	os.CopyFS("build/style/", os.DirFS("style/"))

	fmt.Println("Site generated in /build")
}