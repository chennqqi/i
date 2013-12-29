package i

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

const (
	germanXML   = "i18n/de.xml"
	russianJSON = "i18n/ru.json"
)

func Setup() {
	rand.Seed(int64(os.Getpid()))
}

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

func randomString(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(26) + 97)
	}

	return string(b)
}

func fillStorageWithData(st Storage) {
	const (
		nLocales      = 100
		nScopes       = 10
		nTranslations = 100
		trLen         = 5
	)
	var (
		locales [nLocales]string
		scopes  [nScopes]string
	)
	for i := 0; i < nLocales; i++ {
		locales[i] = randomString(2)
	}
	for i := 0; i < nScopes; i++ {
		scopes[i] = randomString(10)
	}
	for _, locale := range locales {
		for _, scope := range scopes {
			for i := 0; i < nTranslations; i++ {
				st.SetTranslation(randomString(trLen), randomString(trLen), scope, locale)
			}
		}
	}
}

func BenchmarkDefaultStorage(b *testing.B) {
	st := NewDefaultStorage()
	fillStorageWithData(st)
	st.SetTranslation("Пока.", "Bye.", "default", "ru")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Translation("Hello.", "default", "ru")
	}
}
