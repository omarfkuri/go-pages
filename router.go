package router

import (
	"github.com/labstack/echo/v4"
)

type Globals struct {
	Hash string
}

type Site[DB any] struct {
	Dir string
	Pages []*Page[DB]
}

type PageInitFn[DB any] func(e *echo.Echo, db *DB, globals *Globals) error
type PageDataFn[DB any] func(db *DB, globals *Globals) (map[string]any, error)
type PageRenderFn func(ctx echo.Context, data interface{}) error
type PageErrorFn func(ctx echo.Context, err error) error

type Page[DB any] struct {
	Path string
	Dir string
	Info Info
	Init PageInitFn[DB]
	Data PageDataFn[DB]
	Render PageRenderFn
	Error PageErrorFn
}

func NewSite[DB any](
	dir string,) *Site[DB] {

	return &Site[DB]{
		Dir: dir,
	}
}

func (s *Site[DB]) Add(page *Page[DB]) {
	s.Pages = append(s.Pages, page)
}

func NewPage[DB any](
	path string, 
	dir string, 
	info InfoParameters,
	generateFn PageDataFn[DB],
	renderFn PageRenderFn,
	errorFn PageErrorFn,
	init PageInitFn[DB]) *Page[DB] {

	return &Page[DB]{
		Path: path,
		Dir: dir,
		Data: generateFn,
		Render: renderFn,
		Error: errorFn,
		Init: init,
		Info: Info{
			URL: "" + path,
			Title: info.Title,
			Description: info.Description,
			Keywords: info.Keywords,
			Author: info.Author,
			Canonical: info.Canonical,
			Image: info.Image,
			Twitter: info.Twitter,
			OpenGraph: info.OpenGraph,
		},
	}
}

func (sitePtr *Site[DB]) Setup(e *echo.Echo, db *DB, globals *Globals) error {
	var err error = nil

	for _, page := range sitePtr.Pages {
		err = page.Setup(sitePtr.Dir, e, db, globals)
	}

	return err
}

func (pagePtr *Page[DB]) Setup(dir string, e *echo.Echo, db *DB, globals *Globals) error {
	page := *pagePtr
	locales, err := loadLocales(dir + "/" + page.Dir)

	if (err != nil) {
		return err
	}

	e.GET(page.Path, func(ctx echo.Context) error {

		lang := determineLanguage(ctx.Request())

		locale := locales[lang]
		data, err := page.Data(db, globals)

		if err != nil {
			return page.Error(ctx, err)
		}

		data["Globals"] = globals
		data["Locale"] = locale
		data["Info"] = page.Info

		return page.Render(ctx, data)
	})

	if err := page.Init(e, db, globals); err != nil {
		return err
	}

	return nil
}