package bento

import (
	"testing"
)

func TestJustificationExtraSpace(t *testing.T) {
	n := &Box{}
	n.layout.InnerWidth = 20
	n.layout.InnerHeight = 20
	bounds := [][2]int{
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 4},
	}
	tests := []struct {
		tag   string
		hjust Justification
		vjust Justification
		want  [][2]int
	}{
		{
			tag:   "row",
			hjust: Start,
			vjust: Start,
			want: [][2]int{
				{0, 0},
				{2, 0},
				{5, 0},
				{9, 0},
			},
		},
		{
			tag:   "row",
			hjust: End,
			vjust: End,
			want: [][2]int{
				{6, 19},
				{8, 18},
				{11, 17},
				{15, 16},
			},
		},
		{
			tag:   "row",
			hjust: Center,
			vjust: Center,
			want: [][2]int{
				{3, 10},
				{5, 9},
				{8, 9},
				{12, 8},
			},
		},
		{
			tag:   "row",
			hjust: Evenly,
			vjust: Start,
			want: [][2]int{
				{1, 0},
				{4, 0},
				{8, 0},
				{13, 0},
			},
		},
		{
			tag:   "row",
			hjust: Around,
			vjust: Start,
			want: [][2]int{
				{0, 0},
				{3, 0},
				{7, 0},
				{12, 0},
			},
		},
		{
			tag:   "row",
			hjust: Between,
			vjust: Start,
			want: [][2]int{
				{0, 0},
				{4, 0},
				{9, 0},
				{15, 0},
			},
		},
	}
	for _, test := range tests {
		n.Children = nil
		for _, b := range bounds {
			n.Children = append(n.Children, &Box{
				layout: layout{
					OuterWidth:  b[0],
					OuterHeight: b[1],
				},
			})
		}
		n.Tag = test.tag
		n.style = &Style{
			HJust: test.hjust,
			VJust: test.vjust,
		}
		n.justify()
		for i, c := range n.Children {
			if c.X != test.want[i][0] || c.Y != test.want[i][1] {
				t.Fatalf("justification %s, %s, child %d got (%d,%d), want (%d,%d)", test.hjust, test.vjust, i, c.X, c.Y, test.want[i][0], test.want[i][1])
			}
		}
	}
}

func TestJustificationNoExtraSpace(t *testing.T) {
	n := &Box{}
	n.layout.InnerWidth = 14
	n.layout.InnerHeight = 4
	bounds := [][2]int{
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 4},
	}
	tests := []struct {
		tag   string
		hjust Justification
		vjust Justification
		want  [][2]int
	}{
		{
			tag:   "row",
			hjust: Start,
			vjust: Start,
			want: [][2]int{
				{0, 0},
				{2, 0},
				{5, 0},
				{9, 0},
			},
		},
		{
			tag:   "row",
			hjust: End,
			vjust: End,
			want: [][2]int{
				{0, 3},
				{2, 2},
				{5, 1},
				{9, 0},
			},
		},
		{
			tag:   "row",
			hjust: Center,
			vjust: Center,
			want: [][2]int{
				{0, 2},
				{2, 1},
				{5, 1},
				{9, 0},
			},
		},
	}
	for _, test := range tests {
		n.Children = nil
		for _, b := range bounds {
			n.Children = append(n.Children, &Box{
				layout: layout{
					OuterWidth:  b[0],
					OuterHeight: b[1],
				},
			})
		}
		n.Tag = test.tag
		n.style = &Style{
			HJust: test.hjust,
			VJust: test.vjust,
		}
		n.justify()
		for i, c := range n.Children {
			if c.X != test.want[i][0] || c.Y != test.want[i][1] {
				t.Fatalf("justification %s, %s, child %d got (%d,%d), want (%d,%d)", test.hjust, test.vjust, i, c.X, c.Y, test.want[i][0], test.want[i][1])
			}
		}
	}
}
