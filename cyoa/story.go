package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	
	<link href="https://stackpath.bootstrapcdn.com/bootswatch/4.1.1/sketchy/bootstrap.min.css" rel="stylesheet" integrity="sha384-LlPOZK7jvvPEMrrhdqVlEYmX2u/GWKdcq/p7wuVYAUladqNeK7VN1PUZQDmiqlql" crossorigin="anonymous">

	<style>
	body {
		padding-top: 2rem;
		padding-bottom: 2rem;
	}
	</style>

    <title>Choose Your Own Adventure</title>
</head>

<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
        <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </div>
</body>

</html>`

type handler struct {
	s            Story
	t            *template.Template
	pathResolver func(r *http.Request) string
}

func defaultPathResolver(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := defaultPathResolver(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// HandlerOption allow to modify set the Story options
type HandlerOption func(h *handler)

// WithTemplate set the template of a Story
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathResolver set the behavior to convert a path to a chapter
func WithPathResolver(pr func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathResolver = pr
	}
}

// NewHandler creates a new Story HTTP Handler
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathResolver}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

// Story is an adventure to read
type Story map[string]Chapter

// Chapter is an arc of the story
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option is a choice to do in the story
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// JSONStory parse a JSON to a Story
func JSONStory(r io.Reader) (Story, error) {
	var story Story
	d := json.NewDecoder(r)
	err := d.Decode(&story)
	if err != nil {
		return nil, err
	}

	return story, nil
}
