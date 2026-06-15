package main

import (
	"LukasDerBaum/templates"
	"LukasDerBaum/src"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func md_to_HTML(md string) string {
	var buf bytes.Buffer

	md_parser := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	
    err := md_parser.Convert([]byte(md), &buf)
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

func RenderHome(title string) templ.Component {
    return templates.BaseLayout(title, src.HomePage())
}

func RenderLinks(title string) templ.Component {
    return templates.BaseLayout(title, src.Links())
}

func RenderAbout(title string) templ.Component {
    return templates.BaseLayout(title, src.About())
}

func RenderGoals(title string) templ.Component {
    return templates.BaseLayout(title, src.Goals())
}


func build_pages(pages map[string]templ.Component) {
	for name, page := range pages {
		os.Mkdir("build/"+name, 0755)
		f, err := os.Create("build/" + name + "index.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = page.Render(context.Background(), f)
		if err != nil {
			panic(err)
		}
	}
	
}

func main() {
	// md, err := os.ReadFile("content/index.md")
	// if err != nil {
	// 	panic(err)
	// }

	// htmlBody := md_to_HTML(string(md))

	var project_page = map[string]templ.Component{}
	var pages = map[string]templ.Component{
		"":  RenderHome("My Site"),
		"links/":  RenderLinks("Links"),
		"about/":  RenderAbout("About"),
		"goals/":  RenderGoals("Goals"),
		"projects/": buildProjects(&project_page),
	}




	os.RemoveAll("build")
	fmt.Println("Removing build directory...")
	os.Mkdir("build", 0755)
	build_pages(pages)
	build_pages(project_page)
	os.CopyFS("build/style/", os.DirFS("style/"))
	os.CopyFS("build/public/", os.DirFS("public/"))

	fmt.Println("Site generated in /build")
}