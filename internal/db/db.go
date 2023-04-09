package db

import (
	"github.com/divpro/monalizer/internal/parser"
	"golang.org/x/mod/modfile"
	"strings"
	"sync"
)

type whiteList []string

func (wl whiteList) isAllowed(val string) bool {
	if len(wl) == 0 {
		return true
	}
	for _, prefix := range wl {
		if strings.HasPrefix(val, prefix) {
			return true
		}
	}

	return false
}

func Create(ch <-chan parser.Result, prefixes []string) (*DB, <-chan error) {
	db := New()
	wl := whiteList(prefixes)
	errCh := make(chan error)
	go func() {
		var wg sync.WaitGroup
		for item := range ch {
			wg.Add(1)
			go func(item parser.Result) {
				defer wg.Done()
				if item.Source.Err != nil {
					errCh <- item.Source.Err
					return
				}
				name := modfile.ModulePath(item.Source.Result)
				mod := db.Module(name)
				mod.GoVersion = item.Result.Go.Version

				for _, req := range item.Result.Require {
					modName := req.Mod.Path
					if !wl.isAllowed(modName) {
						continue
					}

					reqMod := db.Module(modName)
					mod.Require(reqMod, req.Mod.Version, req.Indirect)
				}
			}(item)
		}
		wg.Wait()
		close(errCh)
	}()

	return db, errCh
}
