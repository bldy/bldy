package module

import (
	"log"
	"net/url"

	"bldy.build/bldy/repository"
	"github.com/grailbio/reflow/syntax"
)

type Module struct {
	syntax.Module
}

func New(u *url.URL) (*Module, error) {
	repo, err := repository.Dial(u)
	if err != nil {
		return nil, err
	}
	r, err := repo.Open(u.Path)
	if err != nil {
		return nil, err
	}

	path := u.Path

	lx := &syntax.Parser{
		File: path,
		Body: r,
		Mode: syntax.ParseModule,
	}
	if err := lx.Parse(); err != nil {
		return nil, err
	}
	save := path
	types, _ := syntax.Stdlib()

	log.Println(types)
	if err := lx.Module.Init(nil, types); err != nil {
		path = save
		return nil, err
	}
	mod := lx.Module
	return &Module{mod}, nil
	/*



		s.path = save
		// Label each toplevel declaration with the module name.
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		base = strings.TrimSuffix(base, ext)
		for _, decl := range lx.Module.Decls {
			decl.Ident = base + "." + decl.Ident
		}
		lx.Module.source = source
		s.modules[path] = lx.Module
		mod = lx.Module
		return &Module{mod}, nil
	*/
}
