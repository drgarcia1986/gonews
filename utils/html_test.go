package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCleanHTML(t *testing.T) {
	var htmls = []struct {
		input    string
		expected string
	}{
		{"<html><p>I&#39;m a <b>gopher</b></p></html>", "I'm a gopher"},
		{"foo<script>alert('hi');</script>", "foo"},
	}

	for _, tt := range htmls {
		actual, err := cleanHTML(tt.input)
		if err != nil {
			t.Fatal("error on clean html: ", err)
		}
		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}
	}
}

func TestGetPreview(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<!doctype html>
			<html lang="en">

			<head>
			  <meta charset="utf-8">
			  <meta http-equiv="x-ua-compatible" content="ie=edge">
			  <meta name="viewport" content="width=device-width, initial-scale=1">

			  <title></title>

			  <link rel="stylesheet" href="css/main.css">
			  <link rel="icon" href="images/favicon.png">
			</head>

			<body>
				<p>Hi from gonews tests</p>

			  <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
			  <script src="js/scripts.js"></script>
			</body>

			</html>
		`))
	}))
	defer ts.Close()

	preview, err := GetPreview(ts.URL)
	if err != nil {
		t.Fatal("error on get preview: ", err)
	}

	actual := strings.TrimSpace(preview)
	expected := "Hi from gonews tests"
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
