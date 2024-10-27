package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "regexp"
    "strings"

    "github.com/gorilla/mux"
)

// ParseMediaWiki converts MediaWiki syntax into HTML
func ParseMediaWiki(text string) string {
    // Headings: `= Heading =`, `== Heading ==`, etc.
    text = parseHeadings(text)
    
    // Bold: `'''bold'''`
    text = parseBold(text)

    // Italic: `''italic''`
    text = parseItalic(text)

    // Internal Links: `[[PageName]]`
    text = parseInternalLinks(text)

    // External Links: `[http://example.com]`
    text = parseExternalLinks(text)

    return text
}

// Parse headings like `= Heading =` and convert to `<h1>Heading</h1>`, etc.
func parseHeadings(text string) string {
    headingRegex := regexp.MustCompile(`(?m)^(={1,6})\s*(.*?)\s*\1$`)
    return headingRegex.ReplaceAllStringFunc(text, func(match string) string {
        parts := headingRegex.FindStringSubmatch(match)
        level := len(parts[1])
        return fmt.Sprintf("<h%d>%s</h%d>", level, parts[2], level)
    })
}

// Parse bold syntax `'''bold'''`
func parseBold(text string) string {
    boldRegex := regexp.MustCompile(`'''(.*?)'''`)
    return boldRegex.ReplaceAllString(text, "<b>$1</b>")
}

// Parse italic syntax `''italic''`
func parseItalic(text string) string {
    italicRegex := regexp.MustCompile(`''(.*?)''`)
    return italicRegex.ReplaceAllString(text, "<i>$1</i>")
}

// Parse internal links `[[PageName]]` to `<a href="/page/PageName">PageName</a>`
func parseInternalLinks(text string) string {
    internalLinkRegex := regexp.MustCompile(`\[\[([^\]]+?)\]\]`)
    return internalLinkRegex.ReplaceAllStringFunc(text, func(match string) string {
        page := internalLinkRegex.FindStringSubmatch(match)[1]
        return fmt.Sprintf(`<a href="/page/%s">%s</a>`, page, page)
    })
}

// Parse external links `[http://example.com]`
func parseExternalLinks(text string) string {
    externalLinkRegex := regexp.MustCompile(`\[(http[^\s]+)\]`)
    return externalLinkRegex.ReplaceAllString(text, `<a href="$1">$1</a>`)
}

// ViewPage renders a MediaWiki-formatted page in HTML
func ViewPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    title := vars["title"]

    // Sample content for demonstration
    content := `= Welcome to the Wiki =
This is a simple page about ''Golang''. Visit the '''Golang Page''' by clicking [[Golang]]. 
To learn more, visit [https://golang.org].`

    // Parse content using MediaWiki parser
    htmlContent := ParseMediaWiki(content)

    // Template for displaying the page
    pageTemplate := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>{{ .Title }}</title>
    </head>
    <body>
        <h1>{{ .Title }}</h1>
        <div>{{ .Content }}</div>
    </body>
    </html>`

    tmpl, err := template.New("page").Parse(pageTemplate)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title   string
        Content template.HTML
    }{
        Title:   title,
        Content: template.HTML(htmlContent),
    }

    tmpl.Execute(w, data)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/page/{title}", ViewPage).Methods("GET")

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
