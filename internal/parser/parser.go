package parser

import (
	"fmt"
	"github.com/divpro/monalizer/internal/source"
	"golang.org/x/mod/modfile"
	"sync"
)

type Result struct {
	Source source.Item
	Result *modfile.File
	Err    error
}

func Parse(srcCh <-chan source.Item) chan Result {
	resCh := make(chan Result, len(srcCh))
	go func() {
		var wg sync.WaitGroup
		for item := range srcCh {
			wg.Add(1)
			go func(item source.Item) {
				defer wg.Done()
				res := Result{
					Source: item,
				}
				if item.Err != nil {
					resCh <- res
					return
				}
				f, err := modfile.Parse("go.mod", item.Result, nil)
				if err != nil {
					res.Err = fmt.Errorf("parse go.mod [ %s ]: %w", item.URL, err)
					resCh <- res
					return
				}
				res.Result = f
				resCh <- res
			}(item)
		}
		wg.Wait()
		close(resCh)
	}()

	return resCh
}
