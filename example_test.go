package i_test

import (
	"fmt"
	"github.com/ainar-g/i"
	"os"
	"strings"
)

func Example() {
	// Load translations from a JSON file.
	i.SetLocale("ru")
	if e := i.LoadJSON("i18n/ru.json"); e != nil {
		panic(e)
	}
	fmt.Println("ru/default:", i.T("Hello."))

	// Or an XML file.
	if e := i.LoadXML("i18n/de.xml"); e != nil {
		panic(e)
	}
	fmt.Println("de/default:", i.T("Hello.", "", "de"))
	// Empty scope means deafult scope, same for locale.

	// Or a string.
	s := `
{
	"en": {
		"cat": {
			"Hello.": "hai thar"
		}
	}
}`
	reader := strings.NewReader(s)
	source := i.NewJSONSource(reader)
	i.LoadFrom(source)

	i.SetLocale("en")
	i.SetScope("cat")
	fmt.Println("en/cat:", i.T("Hello."))
	// Instead of i.Set* we could just write
	//    fmt.Println(i.T("Hello.", "cat", "en"))

	// Output:
	// ru/default: Привет.
	// de/default: Guten Tag.
	// en/cat: hai thar
}

func ExampleTranslator_T() {
	file, e := os.Open("i18n/ru.json")
	if e != nil {
		panic(e)
	}
	defer file.Close()

	// Create the source and storage.
	source := i.NewJSONSource(file)
	storage := i.NewDefaultStorage()

	// Create the custom translator.
	t := i.NewTranslator("default", "ru", source, storage)

	// Load the translations from source.
	if e := t.Load(); e != nil {
		panic(e)
	}

	// Use your custom translator as you would do with the package.
	fmt.Println(t.T("How are you?"))
	// Output: Как дела?
}
