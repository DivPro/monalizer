package render

import (
	"fmt"
	"github.com/divpro/monalizer/internal/db"
	"os"
	"text/template"
)

type Render struct {
	conf *Conf
}

type tplItem struct {
	ID    uint8
	Name  string
	Used  int
	Color string
}

type tplLink struct {
	From  uint8
	To    uint8
	Label string
}

type tplData struct {
	Items []tplItem
	Links []tplLink
}

func New(conf *Conf) *Render {
	return &Render{conf: conf}
}

func (r *Render) Output(store *db.DB) error {
	tpl, err := template.ParseFiles(r.conf.Tpl)
	if err != nil {
		return fmt.Errorf("parse tpl [ %s ]: %w", r.conf.Tpl, err)
	}

	out, err := os.OpenFile(r.conf.Out, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create out file [ %s ]: %w", r.conf.Out, err)
	}
	defer func() {
		_ = out.Close()
	}()

	d := tplData{
		Items: make([]tplItem, 0, len(store.Modules)),
	}
	linksTotal := 0

	var maxUsedBy int
	for _, mod := range store.Modules {
		if len(mod.UsedBy) > maxUsedBy {
			maxUsedBy = len(mod.UsedBy)
		}
	}

	for _, mod := range store.Modules {
		d.Items = append(d.Items, tplItem{
			ID:    mod.ID,
			Name:  mod.Name,
			Used:  len(mod.UsedBy),
			Color: getHeatColor(maxUsedBy, len(mod.UsedBy)),
		})
		linksTotal += len(mod.Requires)
	}
	d.Links = make([]tplLink, 0, linksTotal)
	for _, mod := range store.Modules {
		for _, req := range mod.Requires {
			d.Links = append(d.Links, tplLink{
				From:  req.ID,
				To:    mod.ID,
				Label: req.Version,
			})
		}
	}

	return tpl.Execute(out, d)
}

func getHeatColor(max, cur int) string {
	g := 255 - int(float64(cur)/float64(max)*255)

	return fmt.Sprintf("#%02x%02x%02x", 255, g, 0)
}
