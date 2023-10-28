package runner

import (
	"context"

	"github.com/berquerant/godepdump/display"
	"github.com/berquerant/godepdump/packagesx"
	"github.com/berquerant/godepdump/write"
)

func ListImport(ctx context.Context, patterns []string) error {
	loader := packagesx.New()
	if err := loader.Load(ctx, patterns...); err != nil {
		return err
	}

	type Result struct {
		Src *display.Package `json:"src"`
		Dst *display.Package `json:"dst"`
	}
	for _, p := range loader.List() {
		for _, q := range p.Imports {
			write.JSON(Result{
				Src: display.NewPackage(p),
				Dst: display.NewPackage(q),
			})
		}
	}
	return nil
}
