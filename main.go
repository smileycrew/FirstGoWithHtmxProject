package main

import (
	"example/FirstApi/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"strconv"
	"time"
	// Short answer: used to communication over the network and give development data
	"io"
)

type Template struct {
	templates *template.Template // *template.Template is a type which is used to represent one of more parsed templates
}

// function called Render (based on Template) takes io writer, name, data, context from echo and returns an error
func (template *Template) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	// executes template by writing to writer the template and the data
	return template.templates.ExecuteTemplate(writer, name, data)
}

// newTemplate - create new template instance and parse out the html file
func newTemplate() *Template {
	tmpl := template.Must(template.ParseGlob("*templates/pages/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/partials/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/layouts/*.html"))
	// creates instance and returns a pointer to the template
	return &Template{
		// ParseGlob parses the html files that match the pattern
		// .Must wraps ParseGlub and throws an error if there is an issue parsing the html
		templates: tmpl,
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := models.InitialPageInfo()

	e.Renderer = newTemplate()

	e.Static("/styles", "styles")

	e.GET("/", func(context echo.Context) error {
		return context.Render(200, "index", page)
	})

	e.POST("/contacts", func(context echo.Context) error {
		email := context.FormValue("email")
		name := context.FormValue("name")

		if page.Data.IsEmailTaken(email) {
			formData := models.NewFormData()
			formData.Values["email"] = email
			formData.Values["name"] = name
			formData.Errors["email"] = "Email could not be saved."

			return context.Render(422, "form", formData)
		}

		contact := models.NewContact(email, name)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		context.Render(200, "form", models.NewFormData())

		return context.Render(200, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(context echo.Context) error {
		time.Sleep(3 * time.Second)
		idStr := context.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return context.String(400, "Invalid id.")
		}

		index := page.Data.IndexOf(id)

		if index == -1 {
			return context.String(404, "Contact not found")
		}

		page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)

		return context.NoContent(200)

	})

	e.Logger.Fatal(e.Start(":42069"))
}
