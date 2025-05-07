package repo

import (
	"strings"

	"github.com/cirius-go/generic/slice"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cirius-go/portfolio-server/util"
)

type ParseSortOptions struct {
	separator string
	pipeFns   []slice.PipeFn[string]
}

func (w *ParseSortOptions) Separator(s string) *ParseSortOptions {
	w.separator = s
	return w
}

func (w *ParseSortOptions) PipeFns(fns ...slice.PipeFn[string]) *ParseSortOptions {
	w.pipeFns = append(w.pipeFns, fns...)
	return w
}

func NewWithSortOptions() *ParseSortOptions {
	return &ParseSortOptions{
		separator: ",",
		pipeFns:   make([]slice.PipeFn[string], 0),
	}
}

func QuoteCol(db *gorm.DB, name string) string {
	b := &strings.Builder{}
	db.QuoteTo(b, name)
	return b.String()
}

func ParseSort(db *gorm.DB, sort string, opts ...*ParseSortOptions) clause.OrderBy {
	var (
		opt     = util.IfNull(NewWithSortOptions(), opts...)
		columns = []clause.OrderByColumn{}
	)

	pipes := append([]slice.PipeFn[string]{strings.TrimSpace}, opt.pipeFns...)
	for s := range strings.SplitSeq(sort, opt.separator) {
		s = slice.Pipe(s, pipes...)
		if s == "" {
			continue
		}

		col, found := strings.CutPrefix(s, "-")
		if found {
			columns = append(columns, clause.OrderByColumn{
				Column: clause.Column{Name: QuoteCol(db, col)},
				Desc:   true,
			})
			continue
		}

		col, _ = strings.CutPrefix(s, "+")
		columns = append(columns, clause.OrderByColumn{
			Column: clause.Column{Name: QuoteCol(db, col)},
			Desc:   false,
		})
	}

	return clause.OrderBy{Columns: columns}
}

// page must begin at 1.
func WithPaging(db *gorm.DB, page, perPage int, maxVals ...int) *gorm.DB {
	if page < 1 {
		if len(maxVals) > 0 && maxVals[0] >= 0 {
			return db.Offset(0).Limit(maxVals[0])
		}

		return db.Offset(0).Limit(-1)
	}

	return db.Offset((page - 1) * perPage).Limit(perPage)
}
