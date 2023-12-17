package internals

import (
	"fmt"
	"strings"
	"time"
)

type Author struct {
	name  string
	email string
	time  time.Time
}

func (a *Author) New(name, email string, time time.Time) error {
	a.name = name
	a.email = email
	a.time = time
	return nil
}

func (a *Author) ToString() string {
	return a.name + " <" + a.email + "> " + a.time.Format(time.RFC3339)
}

func ParseAuthor(data string) (*Author, error) {
	d := strings.TrimSpace(data)
	// fmt.Println("inside auth: ", d)
	name := d[:strings.Index(d, "<")]
	email := d[strings.Index(d, "<") : strings.Index(d, ">")-1]
	timeStr := strings.TrimSpace(d[strings.Index(d, ">")+2:])
	// println("name:", name, "email:", email, "time:", timeStr)
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return &Author{
		name:  name,
		email: email,
		time:  t,
	}, nil
}
