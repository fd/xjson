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
		r  = Parse([]byte(js))
		x  Value
		s  string
		b  bool
		ok bool
	)

	x = r.MapIndex("people")
	fmt.Printf("%s => %j\n", x.Selector(), x)

	x = r.MapIndex("people").Index(1).MapIndex("name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.MapIndex("people").Index(1).MapIndex("first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.MapIndex("people").Index(3).MapIndex("first name")
	s, ok = x.MaybeString()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, ok)

	x = r.MapIndex("people").Index(0).MapIndex("name")
	b, ok = x.MaybeBool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, ok)

	x = r.MapIndex("pets").Index(0).MapIndex("name")
	s, ok = x.MaybeString()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, ok)

	// Output:
	// $root.people[1].name => "Hans Spooren"
	// $root.people[1]["first name"] => "Hans"
	// $root.people[3]["first name"] => "" (xjson: index out of range (at: $root.people[3]))
	// $root.people[0].name => false (xjson: string is not a json bool (at: $root.people[0].name))
	// $root.pets[0].name => "" (xjson: key not found (at: $root.pets))
}

func ExampleValue_Path() {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	var (
		r  = Parse([]byte(js))
		x  Value
		s  string
		b  bool
		ok bool
	)

	x = r.Path("people", 1, "name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.Path("people", 1, "first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = r.Path("people", 3, "first name")
	s, ok = x.MaybeString()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, ok)

	x = r.Path("people", 0, "name")
	b, ok = x.MaybeBool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, ok)

	x = r.Path("pets", 0, "name")
	s, ok = x.MaybeString()
	fmt.Printf("%s => %q (%s)\n", x.Selector(), s, ok)

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
		ok     bool
		err    error
	)

	err = json.Unmarshal([]byte(js), &people)
	if err != nil {
		panic(err)
	}

	x = people.People[1].MapIndex("name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = people.People[1].MapIndex("first name")
	fmt.Printf("%s => %q\n", x.Selector(), x.MustString())

	x = people.People[0].MapIndex("name")
	b, ok = x.MaybeBool()
	fmt.Printf("%s => %v (%s)\n", x.Selector(), b, ok)

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

// func ExampleValue_Unwrap() {
//   var js = `
//     {
//       "people": [
//         { "name": "Simon Menke" },
//         { "name": "Hans Spooren", "first name": "Hans" }
//       ]
//     }
//   `

//   type Person map[string]string

//   var (
//     person Person
//     r      = Parse([]byte(js))
//     err    error
//   )

//   err = r.Path("people", 0).Unwrap(&person)
//   fmt.Printf("%#v (err=%s)\n", person, err)

//   err = r.Path("people", 1).Unwrap(&person)
//   fmt.Printf("%#v (err=%s)\n", person, err)

//   // Output:
//   // xjson.Person{"name":"Simon Menke"} (err=%!s(<nil>))
//   // xjson.Person{"name":"Hans Spooren", "first name":"Hans"} (err=%!s(<nil>))
// }
