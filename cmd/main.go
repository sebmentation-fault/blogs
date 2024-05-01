package main

import (
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id = 0

type Contact struct {
	Name, Email string
	Id int
}

func newContact(name, email string) Contact {
	id++
	return Contact{
		Name: name,
		Email: email,
		Id: id,
	}
}

type Contacts = []Contact

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type Data struct {
	Contacts Contacts
}

func (d *Data) indexOf(int int) int {
	for i, contact := range d.Contacts {
		if contact.Id == id {
			return i
		}
	}
	return -1
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Alice", "alice@gmail.com"),
			newContact("Bob", "bob@gmail.com"),
		},
	}
}

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

type Page struct {
	Data Data
	FormData FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		FormData: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")

	page := newPage()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"
			return c.Render(422, "form", formData)
		}

		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		c.Render(200, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func (c echo.Context) error {
		time.Sleep(3 * time.Second)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(400, "Invalid id")
		}

		index := page.Data.indexOf(id)
		if index == -1 {
			return c.String(404, "Contact not found")
		}

		page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)
		return c.NoContent(200)

	})

	e.Logger.Fatal(e.Start(":42069"))
}

