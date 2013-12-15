package i

import (
	"fmt"
	"os"
	"testing"
)

const (
	germanXML   = "i18n/de.xml"
	russianJSON = "i18n/ru.json"
)

func assert(t *testing.T, exp, got, fname string) {
	if got != exp {
		t.Errorf("%s failed, exp: '%s', got '%s'.", fname, exp, got)
	}
}

func TestT(t *testing.T) {
	exp := "nil"
	got := T("nil")
	assert(t, exp, got, "Translator.T")
}

func TestXMLSource(t *testing.T) {
	file, e := os.Open(germanXML)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	source := NewXMLSource(file)
	storage := NewDefaultStorage()
	tr := NewTranslator("default", "de", source, storage)
	tr.Load()

	exp := "Guten Tag."
	got := tr.T("Hello.", "default")
	assert(t, exp, got, "XMLSource")
}

func TestDelete(t *testing.T) {
	file, e := os.Open(russianJSON)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	source := NewJSONSource(file)
	storage := NewDefaultStorage()

	tr := NewTranslator("default", "ru", source, storage)

	// DeleteTranslation
	if e := tr.Load(); e != nil {
		panic(e)
	}

	storage.DeleteTranslation("Hello.", "default", "ru")
	exp := "Hello."
	got := T("Hello.")
	assert(t, exp, got, "DeleteTranslation")

	// DeleteScope
	if e := tr.Load(); e != nil {
		panic(e)
	}

	storage.DeleteScope("default", "ru")
	got = T("Hello.")
	assert(t, exp, got, "DeleteScope")

	// DeleteLocale
	if e := tr.Load(); e != nil {
		panic(e)
	}

	storage.DeleteLocale("ru")
	got = T("Hello.")
	assert(t, exp, got, "DeleteLocale")
}

func TestSetters(t *testing.T) {
	file, e := os.Open(russianJSON)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	source := NewJSONSource(file)
	storage := NewDefaultStorage()

	tr := NewTranslator("nil", "nil", Source(nil), storage)

	tr.SetSource(source)
	tr.SetLocale("ru")
	tr.SetScope("default")

	tr.Load()

	exp := "Привет."
	got := tr.T("Hello.")
	assert(t, exp, got, "Translator.Set*")
}

func ExampleREADME() {
	file, e := os.Open("i18n/ru.json")
	if e != nil {
		panic(e)
	}
	defer file.Close()
	source := NewJSONSource(file)
	storage := NewDefaultStorage()
	mr := NewTranslator("default", "ru", source, storage)
	if e := mr.Load(); e != nil {
		panic(e)
	}
	fmt.Println(mr.T("I PITY THE FOOL WHO DOESN'T MAKE USE OF I18N!"))
	// Output: МНЕ ЖАЛЬ ТОГО ДУРАКА, ЧТО НЕ ИСПОЛЬЗУЕТ ИНТЕРНАЦИОНАЛИЗАЦИЮ!
}
