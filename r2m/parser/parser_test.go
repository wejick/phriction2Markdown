package parser

import (
	"reflect"
	"testing"

	"github.com/wejick/phriction2Markdown/r2m/ast"
)

func Test_isPrefixHeading(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Prefix heading =",
			args: args{
				text: []byte("= it's heading"),
			},
			want: true,
		},
		{
			name: "Prefix heading ==",
			args: args{
				text: []byte("== it's heading"),
			},
			want: true,
		},
		{
			name: "Not prefix heading #",
			args: args{
				text: []byte("# it's not heading"),
			},
			want: false,
		},
		{
			name: "Not prefix heading",
			args: args{
				text: []byte("it's not heading"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrefixHeading(tt.args.text); got != tt.want {
				t.Errorf("isPrefixHeading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_skipCharN(t *testing.T) {
	type args struct {
		text []byte
		char byte
		maxN int
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		{
			name: "one consecutive from one",
			args: args{
				text: []byte("# heading"),
				char: '#',
				maxN: 1,
			},
			wantN: 1,
		},
		{
			name: "two consecutive from three",
			args: args{
				text: []byte("## heading"),
				char: '#',
				maxN: 2,
			},
			wantN: 2,
		},
		{
			name: "four consecutive from five",
			args: args{
				text: []byte("#### heading"),
				char: '#',
				maxN: 5,
			},
			wantN: 4,
		},
		{
			name: "not found",
			args: args{
				text: []byte("heading"),
				char: '#',
				maxN: 2,
			},
			wantN: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := skipCharN(tt.args.text, tt.args.char, tt.args.maxN); gotN != tt.wantN {
				t.Errorf("skipCharN() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_skipChar(t *testing.T) {
	type args struct {
		text []byte
		char byte
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		{
			name: "one consecutive",
			args: args{
				text: []byte("# heading"),
				char: '#',
			},
			wantN: 1,
		},
		{
			name: "two consecutive",
			args: args{
				text: []byte("## heading"),
				char: '#',
			},
			wantN: 2,
		},
		{
			name: "four consecutive",
			args: args{
				text: []byte("#### heading"),
				char: '#',
			},
			wantN: 4,
		},
		{
			name: "not found",
			args: args{
				text: []byte("heading"),
				char: '#',
			},
			wantN: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := skipChar(tt.args.text, tt.args.char); gotN != tt.wantN {
				t.Errorf("skipChar() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_processPrefixHeading(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name       string
		args       args
		wantBlock  ast.Node
		wantOffset int
	}{
		{
			name: "parse header 1",
			args: args{
				text: []byte("= it's header \nit's the next line"),
			},
			wantBlock: &ast.Heading{
				Level: 1,
				Container: ast.Container{
					Content: []byte(" it's header "),
				},
			},
			wantOffset: 14,
		},
		{
			name: "parse header 2",
			args: args{
				text: []byte("== it's header \nit's the next line"),
			},
			wantBlock: &ast.Heading{
				Level: 2,
				Container: ast.Container{
					Content: []byte(" it's header "),
				},
			},
			wantOffset: 15,
		},
		{
			name: "parse no header",
			args: args{
				text: []byte("it's header \nit's the next line"),
			},
			wantBlock:  nil,
			wantOffset: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlock, gotOffset := processPrefixHeading(tt.args.text)
			if !reflect.DeepEqual(gotBlock, tt.wantBlock) {
				t.Errorf("processPrefixHeading() gotBlock = %v, want %v", gotBlock, tt.wantBlock)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("processPrefixHeading() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func Test_skipCharUntil(t *testing.T) {
	type args struct {
		text []byte
		char byte
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		{
			name: "find me found\n",
			args: args{
				text: []byte("012345678\n"),
				char: '\n',
			},
			wantN: 9,
		},
		{
			name: "find me not found\n",
			args: args{
				text: []byte("012345678\n"),
				char: 'x',
			},
			wantN: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := skipCharUntil(tt.args.text, tt.args.char); gotN != tt.wantN {
				t.Errorf("skipCharUntil() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestParser_isHorizontalRule(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
		},
		{
			name: "valid ***",
			args: args{
				text: []byte("***\n"),
			},
			want: true,
		},
		{
			name: "valid * *  * ",
			args: args{
				text: []byte("* *  * \n"),
			},
			want: true,
		},
		{
			name: "valid ---",
			args: args{
				text: []byte("---\n"),
			},
			want: true,
		},
		{
			name: "invalid ---a",
			args: args{
				text: []byte("---a\n"),
			},
			want: false,
		},
		{
			name: "valid ___",
			args: args{
				text: []byte("___\n"),
			},
			want: true,
		},
		{
			name: "invalid _",
			args: args{
				text: []byte("_\n"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			if got := p.isHorizontalRule(tt.args.text); got != tt.want {
				t.Errorf("Parser.isHorizontalRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsHR(b *testing.B) {
	p := &Parser{}

	for n := 0; n < b.N; n++ {
		p.isHorizontalRule([]byte("**************** * \n"))
	}
}

func Test_processHorizontalRule(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name       string
		args       args
		wantBlock  ast.Node
		wantOffset int
	}{
		{
			name:       "empty",
			wantBlock:  &ast.HorizontalRule{},
			wantOffset: 0,
		},
		{
			name: "***\n",
			args: args{
				text: []byte("***\n"),
			},
			wantBlock:  &ast.HorizontalRule{},
			wantOffset: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlock, gotOffset := processHorizontalRule(tt.args.text)
			if !reflect.DeepEqual(gotBlock, tt.wantBlock) {
				t.Errorf("processHorizontalRule() gotBlock = %v, want %v", gotBlock, tt.wantBlock)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("processHorizontalRule() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}
