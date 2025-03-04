/*
1) All the fields in the Movie struct are exported
(i.e. start with a capital letter), which is necessary for them to be visible
to Go’s encoding/json package. Any fields which aren’t exported won’t be included
when encoding a struct to JSON.

2) We’re initializing the Validator instance in the handler and passing it to the
ValidateMovie() function — rather than initializing it in ValidateMovie() and passing it
back as a return value - this is because as our application gets more complex we will
need to call multiple validation helpers from our handlers, rather than just one.
So initializing the Validator in the handler, and then passing it around,
gives us more flexibility.
*/

package data

import (
	"time"

	"greenlight.mazavrbazavr.ru/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	// Use the Runtime type instead of int32. Note that the omitempty directive will
	// still work on this: if the Runtime field has the underlying value 0, then
	// it will be considered empty and omitted -- and the MarshalJSON() method
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres,omitempty"`
	// The version number starts at 1 and will be incremented each time
	// the movie information is updated
	Version int32 `json:"version"`
}

// Validates the Movie struct.
func ValidateMovie(v *validator.Validator, movie *Movie) {
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(
		validator.Unique(movie.Genres),
		"genres",
		"must not contain duplicate values",
	)
}
