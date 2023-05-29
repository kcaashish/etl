package mapper

import (
	"strings"

	"github.com/oarkflow/pkg/evaluate"

	"github.com/oarkflow/etl"
)

type Config struct {
	FieldMaps           map[string]string
	Lookups             map[string][]map[string]any
	LookupFunc          func(data string, key string, value any) map[string]any
	KeepUnmatchedFields bool
}

type Mapper struct {
	cfg *Config
}

func (m *Mapper) Name() string {
	return "mapper"
}

func (m *Mapper) Transform(data etl.Data) error {
	switch data := data.(type) {
	case map[string]any:
		var fields []string
		for f, _ := range data {
			fields = append(fields, f)
		}
		for dest, src := range m.cfg.FieldMaps {
			if strings.HasPrefix(src, "{{") {
				p, _ := evaluate.Parse(src, true)
				pr := evaluate.NewEvalParams(data)
				d, err := p.Eval(pr)
				if err == nil {
					data[dest] = d
				}
			} else if val, ok := data[src]; ok {
				data[dest] = val
			}
		}
		if !m.cfg.KeepUnmatchedFields {
			for k, _ := range data {
				if _, ok := m.cfg.FieldMaps[k]; !ok {
					delete(data, k)
				}
			}
		}
	}
	return nil
}

func New(cfg *Config) *Mapper {
	return &Mapper{
		cfg: cfg,
	}
}

/*
func init() {
	evaluate.AddCustomOperator("lookupIn", lookupIn)
}

func lookupIn(ctx evaluate.EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(3); err != nil {
		return nil, err
	}
	arg1, err := ctx.Arg(0)
	if err != nil {
		return nil, err
	}
	arg2, err := ctx.Arg(1)
	if err != nil {
		return nil, err
	}
	arg3, err := ctx.Arg(2)
	if err != nil {
		return nil, err
	}
	lookup := arg1.(string)
	key := arg2.(string)
	value := arg3.(string)
	return nil, nil
}
*/
