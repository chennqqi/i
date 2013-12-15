# i
i (pronounced as in *bit*), a minimalistic yet flexible internationalization package.

Version 0.1.0. Slightly unstable.

## Basic Usage
```js
// i18n/ru.json
{
	"ru": {
		"default": {
			"Hello.": "Привет."
		}
	}
}
```

```go
package main

import (
	"fmt"
	"github.com/ainar-g/i"
)

func main() {
	i.LoadJSON("i18n/ru.json") // Your file with translations.
	i.SetLocale("ru")
	fmt.Println(i.T("Hello."))
	// Output: "Привет."
}
```

## Advanced Usage
### Scopes and locales
*Locales* are your program's languages and dialects. You can use any format to name your locales, so a Russian locale may be coded as "ru\_RU", "Russian", or just "ru". The latter is recommended.

*Scopes* are a way to resolve translation conflicts. For example, you have a program that has users and groups, and both of them have names, and you want to translate it into Russian. But in Russian the word "name" is translated to "имя" for the animate objects (i.e. people) and "название" for the inanimate objects (i.e. groups). This is where you use a scope:

```go
i.T("Name", "user", "ru")
// Output: Имя
i.T("Name", "group", "ru")
// Output: Название
```

### Translators
Having one scope and one locale for the whole app may be OK for the smaller programs, but with the bigger systems you might get info troubles. That's where the *custom translators* come into play. You can create as many of them as you want, provide them with the source to get the translations from and the storage to keep them, set the locale and the scope, and you are ready to go.

```go
// Open your file...
file, e := os.Open("i18n/ru.json")
if e != nil {
	panic(e)
}
defer file.Close()
// create the source...
source := i.NewJSONSource(file)

// and a storage. We'll take a default one here.
storage := i.NewDefaultStorage()

// Use them to create the custom translator. Let's call him mr.
mr := i.NewTranslator("default", "ru", source, storage)

// Load the translations from source.
if e := mr.Load(); e != nil {
	panic(e)
}

// Use your custom translator as you would do with the package.
fmt.Println(mr.T("I PITY THE FOOL WHO DOESN'T MAKE USE OF I18N!"))
// Output: МНЕ ЖАЛЬ ТОГО ДУРАКА, ЧТО НЕ ИСПОЛЬЗУЕТ ИНТЕРНАЦИОНАЛИЗАЦИЮ!
```

### Storages
*Storage* is where your translations are stored. i provides the ``DefaultStorage`` type, which is suitable for the most i18n needs, but if you want something more, you can create your own source, be it a map of maps or a key-value in-memory storage à la Redis.

### Sources
*Source* is where the package (and the translators) get their translations from. i provides the default JSON and XML sources. You can make sources of your own, be it a database, a CSV file, and what-not.

# GoDoc
See more package documentation at [GoDoc.org](http://godoc.org/github.com/ainar-g/i).

# License
Released under the [MIT License](http://opensource.org/licenses/MIT).
