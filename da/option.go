package da

import "net/http"

type clientOpts struct {
	httpClient *http.Client
}

type Option interface {
	apply(opts *clientOpts)
}

type optionFunc func(opts *clientOpts)

func (f optionFunc) apply(opts *clientOpts) {
	f(opts)
}

func WithHTTPClient(client *http.Client) Option {
	return optionFunc(func(opts *clientOpts) {
		opts.httpClient = client
	})
}
