package main

import (
	"html/template"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// Short answer: used to communication over the network and give development data
	"io"
)

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

// ==============================
var id = 0

type Contact struct {
	Email string
	Id    int
	Name  string
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (data *Data) indexOf(id int) int {
	for i, contact := range data.Contacts {
		if contact.Id == id {
			return i
		}
	}

	return -1
}

// isEmailTaken - returns a boolean if email exists from the Contacts pointer
func (data *Data) isEmailTaken(email string) bool {
	for _, contact := range data.Contacts {
		if contact.Email == email {
			return true
		}
	}

	return false
}

// newContact - creates a new contact struct using the provided email and name
func newContact(email string, name string) Contact {
	id++
	return Contact{
		Email: email,
		Id:    id,
		Name:  name,
	}
}

// initialData - used to initialize the data and provide two contacts when program runs
func initialData() Data {
	return Data{
		Contacts: []Contact{
			newContact("jd@gmail.com", "John"),
			newContact("em@gmail.com", "Ed"),
		},
	}
}

// ==============================
// ==============================
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
	// creates instance and returns a pointer to the template
	return &Template{
		// ParseGlob parses the html files that match the pattern
		// .Must wraps ParseGlub and throws an error if there is an issue parsing the html
		templates: template.Must(template.ParseGlob("*views/*.html")),
	}
}

// ==============================
type Page struct {
	Data Data
	Form FormData
}

// initialPageInfo - returns a page with the initial data and an empty form
func initialPageInfo() Page {
	return Page{
		Data: initialData(),
		Form: newFormData(),
	}
}

// ==============================
func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := initialPageInfo()

	e.Renderer = newTemplate()

	e.GET("/", func(context echo.Context) error {
		return context.Render(200, "index", page)
	})

	e.POST("/contacts", func(context echo.Context) error {
		email := context.FormValue("email")
		name := context.FormValue("name")

		if page.Data.isEmailTaken(email) {
			formData := newFormData()
			formData.Values["email"] = email
			formData.Values["name"] = name
			formData.Errors["email"] = "Email could not be saved."

			return context.Render(422, "form", formData)
		}

		contact := newContact(email, name)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		context.Render(200, "form", newFormData())

		return context.Render(200, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(context echo.Context) error {
		time.Sleep(3 * time.Second)
		idStr := context.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return context.String(400, "Invalid id.")
		}

		index := page.Data.indexOf(id)

		if index == -1 {
			return context.String(404, "Contact not found")
		}

		page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)

		return context.NoContent(200)

	})

	e.Logger.Fatal(e.Start(":42069"))
}
