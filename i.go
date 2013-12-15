/*
i (pronounced as in bit), a minimalistic yet flexible internationalization package.
*/
package i

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
)

var (
	defaultLocale string
	defaultScope  string
	translator    *Translator
)

func init() {
	defaultScope = "default"
	defaultLocale = "en"
	source := Source(nil)
	storage := NewDefaultStorage()

	translator = NewTranslator(defaultScope, defaultLocale, source, storage)
}

// The Storage is used to store and retrieve translations.
//
// Translation is the getter method of the Storage. It returns the translation
// string and the status. If the translation could not be retrieved,
// the status should be false.
//
// NB. Any errors that may be encountered during the translation retrieval
// should be dealt with whithin the Translation method.
//
// SetTranslation sets the translation for the key in the scope of the locale
// and returns any error encountered.
//
// Delete* methods do exactly what the name says. They delete the given
// translationr, scope or locale from the storage. Handle with care.
type Storage interface {
	Translation(key string, scope string, locale string) (translation string, status bool)
	SetTranslation(translation string, key string, scope string, locale string) error
	DeleteTranslation(key string, scope string, locale string) error
	DeleteScope(scope string, locale string) error
	DeleteLocale(locale string) error
}

// The default i's translation storage.
type DefaultStorage struct {
	trMap map[string]map[string]map[string]string
}

// Creates a new DefaultStorage.
func NewDefaultStorage() (storage *DefaultStorage) {
	return &DefaultStorage{trMap: map[string]map[string]map[string]string{}}
}

func (ts *DefaultStorage) Translation(key, scope, locale string) (string, bool) {
	if _, ok := ts.trMap[locale]; !ok {
		return key, ok
	}
	if _, ok := ts.trMap[locale][scope]; !ok {
		return key, ok
	}
	if translation, ok := ts.trMap[locale][scope][key]; !ok {
		return key, ok
	} else {
		return translation, ok
	}
}

func (ts *DefaultStorage) SetTranslation(translation, key, scope, locale string) error {
	if _, ok := ts.trMap[locale]; !ok {
		ts.trMap[locale] = map[string]map[string]string{}
	}
	if _, ok := ts.trMap[locale][scope]; !ok {
		ts.trMap[locale][scope] = map[string]string{}
	}
	ts.trMap[locale][scope][key] = translation

	return nil
}

func (ts *DefaultStorage) DeleteLocale(locale string) error {
	delete(ts.trMap, locale)
	return nil
}

func (ts *DefaultStorage) DeleteScope(scope string, locale string) error {
	delete(ts.trMap[locale], scope)
	return nil
}

func (ts *DefaultStorage) DeleteTranslation(key string, scope string, locale string) error {
	delete(ts.trMap[locale][scope], key)
	return nil
}

// T translates the key string, using scope args[0] and locale args[1],
// if provided, otherwise using the defaults. If the translation, scope or locale
// are not found, the key string is returned.
func T(key string, args ...string) string {
	return translator.T(key, args...)
}

// Loads translations from a JSON file.
// Example of a JSON translation file:
//
//    {
//      "ru": {
//        "default": {
//          "Hello.": "Привет.",
//          "How are you?": "Как дела?",
//        }
//        "preved": {
//          "Hello.": "Превед.",
//          "How are you?": "Кагдила?",
//        }
//      }
//    }
func LoadJSON(filename string) error {
	f, e := os.Open(filename)
	defer f.Close()
	if e != nil {
		return e
	}
	js := JSONSource{f}
	translator.LoadFrom(&js)
	return nil
}

// Loads translations from an XML file.
// Example of an XML translation file:
//
//    <locale name="de">
//      <scope name="default">
//        <translation key="Hello." value="Guten Tag." />
//      </scope>
//      <scope name="bavaria">
//        <translation key="Hello." value="Grüß Gott." />
//      </scope>
//    </locale>
func LoadXML(filename string) error {
	f, e := os.Open(filename)
	defer f.Close()
	if e != nil {
		return e
	}
	xs := XMLSource{f}
	translator.LoadFrom(&xs)
	return nil
}

// Load translations from the source to the package's storage.
func LoadFrom(source Source) error {
	return source.LoadTranslations(&translator.storage)
}

// Sets the translation locale.
func SetLocale(locale string) {
	translator.locale = locale
}

// Sets the translation scope.
func SetScope(scope string) {
	translator.scope = scope
}

// The Translator is the object used to create a custom internationalizer.
// If needed, one can create a number of Translators with different
// locales and scopes rather than load all translations into the package's
// storage and use T with three artuments.
type Translator struct {
	locale  string
	scope   string
	source  Source
	storage Storage
}

// NewTranslator creates a new Translator with the specified scope, locale,
// source and storage.
func NewTranslator(scope, locale string, source Source, storage Storage) *Translator {
	return &Translator{scope: scope, locale: locale, source: source, storage: storage}
}

// T translates the key string, using scope args[0] and locale args[1]
// if provided, or using the Translator's defaults. If the translation,
// scope or locale is not found, the key is returned.
func (t *Translator) T(key string, args ...string) string {
	scope, locale := "", ""

	switch len(args) {
	case 2:
		scope = args[0]
		locale = args[1]
		if scope == "" {
			scope = t.scope
		}
		if locale == "" {
			locale = t.locale
		}
	case 1:
		scope = args[0]
		locale = t.locale
		if scope == "" {
			scope = t.scope
		}
	case 0:
		scope = t.scope
		locale = t.locale
	default:
		panic("Wrong number of arguments for T: need 1 to 3.")
	}
	translation, ok := t.storage.Translation(key, scope, locale)
	if ok {
		return translation
	} else {
		return key
	}
}

// Load translations from the translator's source.
func (t *Translator) Load() error {
	return t.source.LoadTranslations(&t.storage)
}

// Load translations from the source to the translator's storage.
func (t *Translator) LoadFrom(source Source) error {
	return source.LoadTranslations(&t.storage)
}

// Sets the translator's locale.
func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

// Sets the translator's scope.
func (t *Translator) SetScope(scope string) {
	t.scope = scope
}

// Sets the translator's translation source.
func (t *Translator) SetSource(source Source) {
	t.source = source
}

// A Source is where Storage gets the translations from.
//
// LoadTranslations loads translations from the Source to the Storage.
// Returns error.
type Source interface {
	LoadTranslations(*Storage) error
}

type (
	localeXML struct {
		XMLName xml.Name   `xml:"locale"`
		Name    string     `xml:"name,attr"`
		Scopes  []scopeXML `xml:"scope"`
	}

	scopeXML struct {
		XMLName      xml.Name         `xml:"scope"`
		Name         string           `xml:"name,attr"`
		Translations []translationXML `xml:"translation"`
	}

	translationXML struct {
		XMLName xml.Name `xml:"translation"`
		Key     string   `xml:"key,attr"`
		Value   string   `xml:"value,attr"`
	}
)

// XMLSource can be a file or a buffer. Anything that implements io.Reader.
type XMLSource struct {
	reader io.Reader
}

// Creates a new XMLSource.
func NewXMLSource(reader io.Reader) *XMLSource {
	return &XMLSource{reader: reader}
}

func (src *XMLSource) LoadTranslations(ts *Storage) error {
	var locales []localeXML

	translationsFile, e := ioutil.ReadAll(src.reader)
	if e != nil {
		return e
	}
	// If the source is a file, rewind it, so that one could reload the translations.
	if file, ok := src.reader.(*os.File); ok {
		if _, e := file.Seek(0, os.SEEK_SET); e != nil {
			return e
		}
	}

	if e := xml.Unmarshal(translationsFile, &locales); e != nil {
		return e
	}
	for _, locale := range locales {
		for _, scope := range locale.Scopes {
			for _, translation := range scope.Translations {
				key := translation.Key
				value := translation.Value

				if e := (*ts).SetTranslation(value, key, scope.Name, locale.Name); e != nil {
					return e
				}
			}
		}
	}

	return nil
}

// JSONSource can be a file or a buffer. Anything that implements io.Reader.
type JSONSource struct {
	reader io.Reader
}

// Creates a new JSONSource.
func NewJSONSource(reader io.Reader) *JSONSource {
	return &JSONSource{reader: reader}
}

func (src *JSONSource) LoadTranslations(ts *Storage) error {
	translationsFile, e := ioutil.ReadAll(src.reader)
	if e != nil {
		return e
	}
	// If the source is a file, rewind it, so that one could reload the translations.
	if file, ok := src.reader.(*os.File); ok {
		if _, e := file.Seek(0, os.SEEK_SET); e != nil {
			return e
		}
	}

	var locales map[string]map[string]map[string]string
	if e = json.Unmarshal(translationsFile, &locales); e != nil {
		return e
	}

	for locale, scopes := range locales {
		for scope, translations := range scopes {
			for key, value := range translations {
				if e = (*ts).SetTranslation(value, key, scope, locale); e != nil {
					return e
				}
			}
		}
	}

	return nil
}
