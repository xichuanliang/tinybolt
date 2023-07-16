package page

import (
	"testing"
)

func Test_page_typ(t *testing.T) {
	tests := []struct {
		name string
		p    *page
		want string
	}{
		{
			name: "a",
			p:    &page{flags: uint16(1)},
			want: "branch",
		},
		{
			name: "b",
			p:    &page{flags: uint16(18)},
			want: "leaf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.typ(); got != tt.want {
				t.Errorf("page.typ() = %v, want %v", got, tt.want)
			}
		})
	}
}
