package i18n

import (
	"strings"
	"sync"
)

// Locale is a BCP-47 language tag (e.g. "en", "lo", "th", "zh").
type Locale string

const (
	EN Locale = "en"
	LO Locale = "lo" // Lao
	TH Locale = "th" // Thai
	ZH Locale = "zh" // Chinese
)

var (
	mu       sync.RWMutex
	registry = map[Locale]map[string]string{}
)

func init() {
	for locale, messages := range defaultCatalog {
		registry[locale] = messages
	}
}

// Register adds or overwrites translations for a locale.
// Call this at startup to add a new language or override built-in messages.
//
//	i18n.Register(i18n.LO, map[string]string{
//	    "USER_NOT_FOUND": "ບໍ່ພົບຜູ້ໃຊ້",
//	})
func Register(locale Locale, messages map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	if registry[locale] == nil {
		registry[locale] = make(map[string]string)
	}
	for k, v := range messages {
		registry[locale][k] = v
	}
}

// Translate returns the message for code in the given locale.
// Falls back to English if the locale has no entry for code.
// Returns an empty string if the code is unknown in all locales.
func Translate(locale Locale, code string) string {
	mu.RLock()
	defer mu.RUnlock()

	if msgs, ok := registry[locale]; ok {
		if msg, ok := msgs[code]; ok {
			return msg
		}
	}
	// Fallback to English
	if locale != EN {
		if msgs, ok := registry[EN]; ok {
			if msg, ok := msgs[code]; ok {
				return msg
			}
		}
	}
	return ""
}

// FromHeader parses an Accept-Language header and returns the best matching Locale.
// Supports full tags ("lo-LA") and short tags ("lo"). Falls back to EN.
//
//	locale := i18n.FromHeader(c.Get("Accept-Language"))
func FromHeader(header string) Locale {
	if header == "" {
		return EN
	}
	// Take the first language tag (before comma), strip quality weight (before semicolon)
	first := strings.TrimSpace(strings.Split(strings.Split(header, ",")[0], ";")[0])
	full := Locale(strings.ToLower(first))
	prefix := Locale(strings.ToLower(strings.Split(first, "-")[0]))

	mu.RLock()
	defer mu.RUnlock()

	if _, ok := registry[full]; ok {
		return full
	}
	if _, ok := registry[prefix]; ok {
		return prefix
	}
	return EN
}

// Supported returns all registered locale codes.
func Supported() []Locale {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]Locale, 0, len(registry))
	for l := range registry {
		list = append(list, l)
	}
	return list
}
