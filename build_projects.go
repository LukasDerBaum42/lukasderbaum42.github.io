package main
import (
	"github.com/a-h/templ"
	"LukasDerBaum/templates"
	"LukasDerBaum/src"
	"os"
	"strings"
	"github.com/BurntSushi/toml"
	// "fmt"
)


type FrontMatter struct {
    Title       string
    Description string
}


func splitFrontMatter(md string) (FrontMatter, string) {
    parts := strings.SplitN(md, "+++", 3)

    if len(parts) < 3 {
        return FrontMatter{}, md
    }

    rawTOML := parts[1]
    content := parts[2]

    var fm FrontMatter
    toml.Decode(rawTOML, &fm)

    return fm, content
}


func buildProjectPage(title string, body string) templ.Component {
	return templates.BaseLayout(title, src.ProjectPage(body))
}



func buildProjects(pages *map[string]templ.Component) templ.Component {
	files, err := os.ReadDir("content/projects")
	if err != nil {
		panic(err)
	}
	
	var project_list []templ.Component

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		md_bytes, err := os.ReadFile("content/projects/" + file.Name())
		
		fm, content := splitFrontMatter(string(md_bytes))
		if err != nil {
			panic(err)
		}
		parsedContent := md_to_HTML(string(content))
		file_name := strings.TrimSuffix(file.Name(), ".md")
		path := "projects/" + file_name + "/"
		(*pages)[path] = buildProjectPage(fm.Title, parsedContent)
		tile := src.ProjectTile(fm.Title,fm.Description,file_name)
		project_list = append(project_list, tile)
	}

	return templates.BaseLayout("Projects", src.Project(project_list))
}
