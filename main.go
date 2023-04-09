package main

import (
	"context"
	"flag"
	"github.com/divpro/monalizer/internal/app"
	"github.com/divpro/monalizer/internal/db"
	"github.com/divpro/monalizer/internal/parser"
	"github.com/divpro/monalizer/internal/render"
	"github.com/divpro/monalizer/internal/source"
	"log"
	"time"
)

var (
	confFile = flag.String("conf", "cfg/conf.yaml", "")
)

func main() {
	flag.Parse()
	conf, err := app.Load(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	src := source.New(conf.Source)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	store, doneCh, errCh := db.Create(parser.Parse(src.Load(ctx)), conf.Whitelist)
	go func() {
		for err := range errCh {
			log.Println(err)
		}
	}()
	<-doneCh
	out := render.New(conf.Render)
	if err := out.Output(store); err != nil {
		log.Println(err)
	}
}
