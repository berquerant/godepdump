package chanx_test

import (
	"testing"

	"github.com/berquerant/godepdump/chanx"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var (
		NilFilter        chanx.Filter[int]
		PositiveFilter   chanx.Filter[int] = func(x int) bool { return x > 0 }
		LessThan10Filter chanx.Filter[int] = func(x int) bool { return x < 10 }
	)

	for _, tc := range []struct {
		name   string
		filter chanx.Filter[int]
		arg    int
		want   bool
	}{
		{
			name:   "nil filter return true",
			filter: NilFilter,
			arg:    100,
			want:   true,
		},
		{
			name:   "positive filter return true",
			filter: PositiveFilter,
			arg:    100,
			want:   true,
		},
		{
			name:   "nil and positive filter return false",
			filter: NilFilter.And(PositiveFilter),
			arg:    0,
			want:   false,
		},
		{
			name:   "between filter return true",
			filter: PositiveFilter.And(LessThan10Filter),
			arg:    4,
			want:   true,
		},
		{
			name:   "between filter return false",
			filter: LessThan10Filter.And(PositiveFilter),
			arg:    100,
			want:   false,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.filter.Call(tc.arg))
		})
	}
}
