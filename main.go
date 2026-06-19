package main

import (
	"LukasDerBaum/src"
	"LukasDerBaum/templates"
	"LukasDerBaum/src/i18n-de"
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

func RenderHome(title, path, prefix string) templ.Component {
    return templates.BaseLayout(title,path,prefix, src.HomePage())
}

func RenderLinks(title, path, prefix string) templ.Component {
    return templates.BaseLayout(title,path,prefix, src.Links())
}

func RenderAbout(title, path, prefix string) templ.Component {
    return templates.BaseLayout(title, path,prefix, src.About())
}

func RenderGoals(title, path, prefix string) templ.Component {
    return templates.BaseLayout(title, path,prefix, src.Goals())
}

func build_german(pages, project_page *map[string]templ.Component){
	os.MkdirAll("build/de/", 0755)
	prefix := "/de"
	(*pages)["de/"] = templates.BaseLayout("Meine Seite", "de/", prefix, src_i18n_de.HomePage())
	(*pages)["de/projects/"] = buildProjects("/de","projects/", project_page)
	(*pages)["de/about/"] = templates.BaseLayout("About me", "de/about/", prefix, src.About())
	(*pages)["de/goals/"] = templates.BaseLayout("Goals", "de/goals/", prefix, src.Goals())
	(*pages)["de/links/"] = templates.BaseLayout("Links", "de/links/", prefix, src.Links())
}


func build_pages(pages map[string]templ.Component, dev bool) {
	for name, page := range pages {
		os.MkdirAll("build/"+name, 0755)
		f, err := devCreate("build/" + name + "index.html", dev)
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
	dev := len(os.Args) > 1 && os.Args[1] == "--dev"

	if !dev {
		os.RemoveAll("build")
		fmt.Println("Removing build directory...")
	} else {
		fmt.Println("Dev mode: keeping existing build directory")
	}
	os.MkdirAll("build", 0755)
	prefix := ""
	var project_page = map[string]templ.Component{}
	var pages = map[string]templ.Component{
		"":  RenderHome("My Site","", prefix),
		"links/":  RenderLinks("Links","links/", prefix),
		"about/":  RenderAbout("About","about/",prefix),
		"goals/":  RenderGoals("Goals","goals/",prefix),
		"projects/": buildProjects(prefix, "projects/",&project_page),
	}

	build_german(&pages, &project_page)


	build_pages(pages, dev)
	build_pages(project_page, dev)
	devCopyFS("build/style/", os.DirFS("style/"), dev)
	devCopyFS("build/public/", os.DirFS("public/"), dev)
	// compilTS(dev)
	

	fmt.Println("Site generated in /build")
}