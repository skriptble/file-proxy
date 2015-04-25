package froxy

import "io"

type proxy struct {
	sources map[string]FileSource
}

func NewProxy() ProxyBuilder {
	p := new(proxy)
	p.sources = make(map[string]FileSource)
	return p
}

func (p *proxy) RetrieveFile(name string, src string) (io.ReadCloser, error) {
	fs, exists := p.sources[src]
	if !exists {
		return nil, NotFound
	}
	rc, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func (p *proxy) AddFileSource(fs FileSource, src string) error {
	p.sources[src] = fs
	return nil
}
