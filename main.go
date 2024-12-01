package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

// BlogPost represents a single blog post
type BlogPost struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Date        string `json:"date"`
	Text        string `json:"text"`
	IsPublished bool   `json:"is_published"`
}

// BlogData represents the entire JSON structure
type BlogData struct {
	Posts []BlogPost `json:"posts"`
}

func main() {
	// Static file server setup
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Updated routes
	http.HandleFunc("/blog/", handleBlog)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Start the server on port 8080
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}

func handleBlog(w http.ResponseWriter, r *http.Request) {
	// Read and parse JSON
	data, err := ioutil.ReadFile("static/data/blog-posts.json")
	if err != nil {
		http.Error(w, "Error reading blog posts", http.StatusInternalServerError)
		return
	}

	var blogData BlogData
	if err := json.Unmarshal(data, &blogData); err != nil {
		http.Error(w, "Error parsing blog posts", http.StatusInternalServerError)
		return
	}

	// Check if we're requesting a specific post
	path := strings.TrimPrefix(r.URL.Path, "/blog/")
	if path != "" && path != "blog" {
		// Show individual post
		for _, post := range blogData.Posts {
			if post.Slug == path && post.IsPublished {
				renderPost(w, post)
				return
			}
		}
		http.NotFound(w, r)
		return
	}

	// Show list of posts
	renderBlogList(w, blogData)
}

func renderBlogList(w http.ResponseWriter, blogData BlogData) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Blog Posts</title>
		<style>
			.card {
				border: 1px solid #ddd;
				padding: 15px;
				margin: 10px;
				border-radius: 5px;
				max-width: 300px;
				display: inline-block;
			}
			.date {
				color: #666;
				font-size: 0.9em;
			}
			a {
				text-decoration: none;
				color: inherit;
			}
		</style>
	</head>
	<body>
		<h1>Blog Posts</h1>
		{{range .Posts}}
			{{if .IsPublished}}
			<a href="/blog/{{.Slug}}">
				<div class="card">
					<h2>{{.Title}}</h2>
					<div class="date">{{.Date}}</div>
				</div>
			</a>
			{{end}}
		{{end}}
	</body>
	</html>
	`

	t, err := template.New("bloglist").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error creating template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, blogData); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func renderPost(w http.ResponseWriter, post BlogPost) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}}</title>
		<style>
			.container {
				max-width: 800px;
				margin: 0 auto;
				padding: 20px;
			}
			.date {
				color: #666;
				font-size: 0.9em;
				margin-bottom: 20px;
			}
			.back-link {
				display: block;
				margin-bottom: 20px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<a href="/blog" class="back-link">← Back to all posts</a>
			<h1>{{.Title}}</h1>
			<div class="date">{{.Date}}</div>
			<div class="content">
				{{.Text}}
			</div>
		</div>
	</body>
	</html>
	`

	t, err := template.New("post").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error creating template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, post); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

//go run main.go
