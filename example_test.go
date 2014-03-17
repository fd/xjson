package xjson

import (
	"encoding/json"
	"fmt"
	"os"
)

func Example() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	var (
		r   = Parse([]byte(js))
		x   Value
		s   string
		b   bool
		err error
	)

	x = r.Get("people").GetIndex(1).Get("name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.Get("people").GetIndex(1).Get("first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.Get("people").GetIndex(3).Get("first name")
	s, err = x.String()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, err)

	x = r.Get("people").GetIndex(0).Get("name")
	b, err = x.Bool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, err)

	x = r.Get("pets").GetIndex(0).Get("name")
	s, err = x.String()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, err)

	// Output:
	// $root.people[1].name => "Hans Spooren"
	// $root.people[1]["first name"] => "Hans"
	// $root.people[3]["first name"] => "" (xjson: index out of range (at: $root.people[3]))
	// $root.people[0].name => false (xjson: string is not a json bool (at: $root.people[0].name))
	// $root.pets[0].name => "" (xjson: key not found (at: $root.pets))
}

func ExampleValue_GetPath() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	var (
		r   = Parse([]byte(js))
		x   Value
		s   string
		b   bool
		err error
	)

	x = r.GetPath("people", 1, "name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.GetPath("people", 1, "first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.GetPath("people", 3, "first name")
	s, err = x.String()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, err)

	x = r.GetPath("people", 0, "name")
	b, err = x.Bool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, err)

	x = r.GetPath("pets", 0, "name")
	s, err = x.String()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, err)

	// Output:
	// $root.people[1].name => "Hans Spooren"
	// $root.people[1]["first name"] => "Hans"
	// $root.people[3]["first name"] => "" (xjson: index out of range (at: $root.people[3]))
	// $root.people[0].name => false (xjson: string is not a json bool (at: $root.people[0].name))
	// $root.pets[0].name => "" (xjson: key not found (at: $root.pets))
}

func ExampleValue_UnmarshalJSON() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	type people_t struct {
		People []Value
	}

	var (
		people people_t
		x      Value
		b      bool
		err    error
	)

	err = json.Unmarshal([]byte(js), &people)
	if err != nil {
		panic(err)
	}

	x = people.People[1].Get("name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = people.People[1].Get("first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = people.People[0].Get("name")
	b, err = x.Bool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, err)

	// Output:
	// $root.name => "Hans Spooren"
	// $root["first name"] => "Hans"
	// $root.name => false (xjson: string is not a json bool (at: $root.name))
}

func ExampleValue_MarshalJSON() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	type people_t struct {
		People []Value
	}

	var (
		people people_t
		err    error
	)

	err = json.Unmarshal([]byte(js), &people)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(people)
	if err != nil {
		panic(err)
	}

	// Output:
	// {"People":[{"name":"Simon Menke"},{"first name":"Hans","name":"Hans Spooren"}]}
}

func ExampleValue_Unwrap() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	type Person map[string]string

	var (
		person Person
		r      = Parse([]byte(js))
		err    error
	)

	err = r.GetPath("people", 0).Unwrap(&person)
	fmt.Printf("%#v (err=%s)\n", person, err)

	err = r.GetPath("people", 1).Unwrap(&person)
	fmt.Printf("%#v (err=%s)\n", person, err)

	// Output:
	// xjson.Person{"name":"Simon Menke"} (err=%!s(<nil>))
	// xjson.Person{"name":"Hans Spooren", "first name":"Hans"} (err=%!s(<nil>))
}
