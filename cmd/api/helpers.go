package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "mooveit-backend.mooveit.com/internal/jsonlog"
	"mooveit-backend.mooveit.com/internal/validator"
)

// Define an envelope type
type envelope map[string]any

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(request *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We can use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(request.Context())

	// We can then use the ByName() method to get the value of the "id" parameter from
	// the slice. In our project all movies will have a unique positive integer ID, but
	// the value returned by ByName() is always a string. So we try to convert it to a
	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
	// or is less than 1, we know the ID is invalid so we use the http.NotFound()
	// function to return a 404 Not Found response.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(writer http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	// Or use the json.MarshalIndent() function so that whitespace is added to the encoded
	// JSON. json.MarshalIndent(data, "", "\t") - here we use no line prefix ("") and tab indents ("\t") for each element.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		writer.Header()[key] = value
	}

	// Set the "Content-Type: application/json" header on the response. If you forget to
	// this, Go will default to sending a "Content-Type: text/plain; charset=utf-8"
	// header instead.
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(js)

	return nil
}

// serverErrorResponse sends a JSON-formatted error message to the client with the given
// status code, and logs the error using our custom logger at the ERROR level.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.ErrorWithProperties(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})

	message := "The server encountered a problem and could not process your request"
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and exit. We don't want to send a response after this point
	// as we will already have sent the HTTP status code to the client.
	err = app.writeJSON(w, http.StatusInternalServerError, env, nil)
	if err != nil {
		log.Error("%s", err)
	}
}

// For a public-facing API, the error messages themselves aren't ideal.
// Some are too detailed and expose information about the underlying
// API implementation. Others aren’t descriptive enough (like "EOF"),
// and some are just plain confusing and difficult to understand.
// There isn’t consistency in the formatting or language used either.
//
//	{
//	   "error": "invalid character '}' looking for beginning of object key string"
//	}
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(destination)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Use the errors.As() function to check whether the error has the type
		// *http.MaxBytesError. If it does, then it means the request body exceeded our
		// size limit of 1MB and we return a clear error message.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic,
		// rather than returning an error to our handler.
		// Bringing this back to our readJSON() helper, if we get a
		// json.InvalidUnmarshalError at runtime it’s because we as the developers
		// have passed an unsupported value to Decode(). This is firmly an unexpected error
		// which we shouldn’t see under normal operation, and is something that should be
		// picked up in development and tests long before deployment.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func processImageData(data any) error {
	// Type assert the data to access the image field
	imageData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid data structure")
	}

	// Check if there's an "image" field in the JSON
	imageStr, ok := imageData["image"].(string)
	if !ok {
		return errors.New("image field not found or not a string")
	}

	// Decode the base64 image data
	imgData, err := base64.StdEncoding.DecodeString(imageStr)
	if err != nil {
		return fmt.Errorf("error decoding base64 image: %v", err)
	}

	// You can now process the image data as needed
	// For example, you might want to validate the image format, resize it, etc.
	// Here we'll just check if it's a valid image
	_, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("invalid image data: %v", err)
	}

	// Update the original data with the processed image information
	imageData["imageFormat"] = format
	imageData["imageSize"] = len(imgData)

	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	str := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if str == "" {
		return defaultValue
	}

	// Otherwise return the string.
	return str
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}

	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(queryString url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	str := queryString.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if str == "" {
		return defaultValue
	}

	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	intValue, err := strconv.Atoi(str)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	// Otherwise, return the converted integer value.
	return intValue
}

// The background() helper accepts an arbitrary function as a parameter.
func (app *application) background(fn func()) {
	// Increment the WaitGroup counter.
	app.wg.Add(1)

	// Launch a background goroutine.
	go func() {
		// Run a deferred function which uses recover() to catch any panic, and log an
		// error message instead of terminating the application.
		defer func() {
			// Use defer to decrement the WaitGroup counter before the goroutine returns.
			defer app.wg.Done()

			if err := recover(); err != nil {
				log.Error("%s", err)
			}
		}()

		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
