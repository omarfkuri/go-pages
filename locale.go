package router

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

type Locales map[string]map[string]any

func loadLocale(dir string, lang string) (*map[string]any, error) {
    filePath := filepath.Join(dir, "locales", lang+".json")

    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var locale map[string]any
    if err := json.NewDecoder(file).Decode(&locale); err != nil {
        return nil, err
    }

    return &locale, nil
}

func loadLocales(dir string) (Locales, error) {
    locales := map[string]map[string]any{}

    en, err := loadLocale(dir, "en");
    if err != nil {
        return locales, err
    }
    locales["en"] = *en;

    es, err := loadLocale(dir, "es");
    if err != nil {
        return locales, err
    }
    locales["es"] = *es;

    return locales, nil
}

func determineLanguage(r *http.Request) string {
    cookie, err := r.Cookie("lang")
    if err == nil {
        return cookie.Value
    }

    acceptLang := r.Header.Get("Accept-Language")
    if acceptLang != "" {
        languages := strings.Split(acceptLang, ",")
        if len(languages) > 0 {
            language := strings.Split(languages[0], ";")[0]
            return strings.Split(language, "-")[0]
        }
    }

    return "en"
}