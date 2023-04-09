package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Source struct {
	conf *Conf
}

type Item struct {
	URL    string
	Result []byte
	Err    error
}

func New(conf *Conf) *Source {
	return &Source{conf: conf}
}

func (s *Source) Load(ctx context.Context) <-chan Item {
	resCh := make(chan Item, len(s.conf.URLs))
	go func() {
		headers := http.Header{}
		for h, val := range s.conf.Headers {
			headers.Set(h, val)
		}

		client := http.Client{
			Timeout: time.Second * 10,
		}
		var wg sync.WaitGroup
		wg.Add(len(s.conf.URLs))
		for _, url := range s.conf.URLs {
			go func(url string) {
				defer wg.Done()
				res := Item{
					URL: url,
				}
				defer func() {
					resCh <- res
				}()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					res.Err = fmt.Errorf("create new request [ %s ]: %w", url, err)
					return
				}
				req.Header = headers
				resp, err := client.Do(req)
				if err != nil {
					res.Err = fmt.Errorf("get response [ %s ]: %w", url, err)
					return
				}
				defer func() {
					if err := resp.Body.Close(); err != nil {
						res.Err = fmt.Errorf("close response body [ %s ]: %w", url, err)
					}
				}()
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					res.Err = fmt.Errorf("read response [ %s ]: %w", url, err)
					return
				}
				res.Result = b
			}(url)
		}
		wg.Wait()
		close(resCh)
	}()

	return resCh
}
