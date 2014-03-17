package xjson

import "testing"

func TestParse(t *testing.T) {
	var js = `
    {
      "people": [
        { "name": "Simon Menke" },
        { "name": "Hans Spooren", "first name": "Hans" }
      ]
    }
  `

	s := new_scanner([]byte(js))
	v, err := s.scan_value()
	t.Logf("v=%+v, err=%s\n", v, err)
}
