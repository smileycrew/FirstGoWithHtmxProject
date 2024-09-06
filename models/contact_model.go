package models

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

// indexOf returns boolean if provided id exists in the []Contacts
func (data *Data) IndexOf(id int) int {
	for i, contact := range data.Contacts {
		if contact.Id == id {
			return i
		}
	}

	return -1
}

// isEmailTaken - returns a boolean if email exists from the Contacts pointer
func (data *Data) IsEmailTaken(email string) bool {
	for _, contact := range data.Contacts {
		if contact.Email == email {
			return true
		}
	}

	return false
}

// newContact - creates a new contact struct using the provided email and name
func NewContact(email string, name string) Contact {
	id++
	return Contact{
		Email: email,
		Id:    id,
		Name:  name,
	}
}

// newData - used to initialize the data and provide two contacts when program runs
func NewData() Data {
	return Data{
		Contacts: []Contact{
			NewContact("jd@gmail.com", "John"),
			NewContact("em@gmail.com", "Ed"),
		},
	}
}

type Page struct {
	Data Data
	Form FormData
}

// initialPageInfo - returns a page with the initial data and an empty form
func InitialPageInfo() Page {
	return Page{
		Data: NewData(),
		Form: NewFormData(),
	}
}
