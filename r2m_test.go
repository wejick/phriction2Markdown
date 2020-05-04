package main

import (
	"testing"
)

func Test_convertHeading(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name            string
		args            args
		wantRevisedLine string
	}{
		{
			name: "H1",
			args: args{
				line: "= baju baru =",
			},
			wantRevisedLine: "#  baju baru ",
		},
		{
			name: "H2",
			args: args{
				line: "== baju baru ==",
			},
			wantRevisedLine: "##  baju baru ",
		},
		{
			name: "H3",
			args: args{
				line: "=== baju baru ===",
			},
			wantRevisedLine: "###  baju baru ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRevisedLine := convertHeading(tt.args.line); gotRevisedLine != tt.wantRevisedLine {
				t.Errorf("fixHeading() = %v, want %v", gotRevisedLine, tt.wantRevisedLine)
			}
		})
	}
}

func Test_convertOrderedList(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name            string
		args            args
		wantRevisedLine string
	}{
		{
			name: " ordered 1",
			args: args{
				line: "# baju baru",
			},
			wantRevisedLine: "1. baju baru",
		},
		{
			name:            "not ordered",
			args:            args{line: "# baju baru #"},
			wantRevisedLine: "# baju baru #",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRevisedLine := convertOrderedList(tt.args.line); gotRevisedLine != tt.wantRevisedLine {
				t.Errorf("convertOrderedList() = %v, want %v", gotRevisedLine, tt.wantRevisedLine)
			}
		})
	}
}
