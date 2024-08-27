package handlers

import (
	"encoding/json"
	"example/FirstApi/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	headerKey   = "Content-Type"
	headerValue = "application/json"
)

var Items []models.Item

func GetItems(writer http.ResponseWriter, request *http.Request) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// encodes the items into JSON and writes the data directly into the writer
	json.NewEncoder(writer).Encode(Items)
}

func getItem(writer http.ResponseWriter, request *http.Request) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// gets the request url and removes the provided prefix
	idStr := strings.TrimPrefix(request.URL.Path, "/items/")

	// try parses the provided string to integer
	id, err := strconv.Atoi(idStr)

	// if parsing was unsuccessful
	if err != nil {
		// http error abstracts this by providing default write, error message and status code
		http.Error(writer, "Invalid item ID", http.StatusBadRequest)

		return
	}

	// _ ignores the index value provided by the range
	for _, item := range Items {
		// we found a match!
		if item.ID == id {
			// encodes the items into JSON and writes the data directly into the writer
			json.NewEncoder(writer).Encode(item)

			// this return breaks out of the loop and the function
			return
		}
	}

	// if everything else fails the default error is sent to the client
	http.Error(writer, "Item not found", http.StatusNotFound)
}

func createItem(writer http.ResponseWriter, request *http.Request) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// initialize an empty instance of an item
	var newItem models.Item

	// try parses the request body sent by the client and assigns it to the referenced variable provided
	err := json.NewDecoder(request.Body).Decode(&newItem)

	// if parsing failed send error to client and return out of the function
	if err != nil {
		http.Error(writer, "Invalid input", http.StatusBadRequest)

		return
	}

	// get the current date time
	now := time.Now()

	// convert the now variable to a string
	nowToStr := now.Format("20060102150405")

	// convert the noToStr variable to an int and assign it as the id
	newItem.ID, err = strconv.Atoi(nowToStr)

	// if there was an error parsing send http error to client
	if err != nil {
		http.Error(writer, "Something went wrong", http.StatusBadRequest)
	}

	// appends the new item to the array
	Items = append(Items, newItem)

	// add the status created header to the writer
	writer.WriteHeader(http.StatusCreated)

	// adds the new item to the header body
	json.NewEncoder(writer).Encode(newItem)
}

func updateItem(writer http.ResponseWriter, request *http.Request) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// gets the request url and removes the provided prefix
	idStr := strings.TrimPrefix(request.URL.Path, "/items/")

	// try parses the provided string to integer
	id, err := strconv.Atoi(idStr)

	// if parsing was unsuccessful
	if err != nil {
		// http error abstracts this by providing default write, error message and status code
		http.Error(writer, "Invalid input", http.StatusBadRequest)

		return
	}

	// initialize a new instance of an item
	var updatedItem models.Item

	// converts the request body to json and try to assign it to the reference variable
	if err := json.NewDecoder(request.Body).Decode(&updatedItem); err != nil {
		http.Error(writer, "Invalid input", http.StatusBadRequest)
	}

	// iterates items and provides index and item
	for i, item := range Items {
		if item.ID == id {
			// if the item exists by matching ids, assigns new values
			Items[i].Name = updatedItem.Name
			Items[i].Price = updatedItem.Price

			// writer is given updated item
			json.NewEncoder(writer).Encode(Items[i])

			// breaks out of the function
			return
		}
	}

	// default will return error to the client
	http.Error(writer, "Item not found", http.StatusNotFound)
}

func deleteItem(writer http.ResponseWriter, request *http.Request) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// gets the request url and removes the provided prefix
	idStr := strings.TrimPrefix(request.URL.Path, "/items/")

	// try parses the provided string to integer
	id, err := strconv.Atoi(idStr)

	// if parsing was unsuccessful
	if err != nil {
		// http error abstracts this by providing default write, error message and status code
		http.Error(writer, "Invalid item id", http.StatusBadRequest)

		return
	}

	// iterates items and provides index and item
	for i, item := range Items {
		// if the item exists by matching ids, items becomes new array without the selected item
		if item.ID == id {
			// items becomes everything before i and everything after i
			Items = append(Items[:i], Items[i+1:]...) // the : creates slices of the values sort of like a reference

			// adds no content to the writer
			writer.WriteHeader(http.StatusNoContent)

			return
		}
	}

	// default error will be sent to the client
	http.Error(writer, "Item not found", http.StatusNotFound)
}

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	// depending the request method it will route to the correct function
	switch request.Method {
	case http.MethodGet:
		getItem(writer, request)
	case http.MethodPost:
		createItem(writer, request)
	case http.MethodPut:
		updateItem(writer, request)
	case http.MethodDelete:
		deleteItem(writer, request)
	default:
		http.Error(writer, "Method not alowed", http.StatusMethodNotAllowed)
	}
}
