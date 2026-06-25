package issue

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []int
		wantErr bool
	}{
		{"single issue", []string{"#123"}, []int{123}, false},
		{"bare number", []string{"123"}, []int{123}, false},
		{"comma separated", []string{"#123,#456"}, []int{123, 456}, false},
		{"and separated", []string{"#123 and #456"}, []int{123, 456}, false},
		{"range", []string{"#10-#15"}, []int{10, 11, 12, 13, 14, 15}, false},
		{"mixed", []string{"#1 and #5-#7"}, []int{5, 6, 7, 1}, false},
		{"deduplicates", []string{"#5,#5,#5"}, []int{5}, false},
		{"no issues", []string{"hello"}, nil, true},
		{"multiple args", []string{"#1", "#2"}, []int{1, 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
