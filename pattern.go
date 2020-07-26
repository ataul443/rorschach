package main

func MatchWildcard(pattern, text string) bool {
	nt := len(text)
	np := len(pattern)

	dp := NewMatrix(nt+1, np+1)

	dp.Set(0, 0, true)

	for j := 1; j <= np; j++ {
		if pattern[j-1] == '*' {
			dp.Set(0, j, dp.Get(0, j-1))
		}
	}

	for i := 1; i <= nt; i++ {
		for j := 1; j <= np; j++ {
			if text[i-1] == pattern[j-1] || pattern[j-1] == '?' {
				dp.Set(i, j, dp.Get(i-1, j-1))

			} else if pattern[j-1] == '*' {
				dp.Set(i, j, dp.Get(i-1, j) || dp.Get(i, j-1))

			}
		}
	}
	return dp.Get(nt, np)
}

type Matrix struct {
	data []bool
	rows int
	cols int
}

func NewMatrix(rowSize, colSize int) *Matrix {
	return &Matrix{
		data: make([]bool, rowSize*colSize),
		rows: rowSize,
		cols: colSize,
	}
}

func (m Matrix) index(rx, cy int) int {
	return rx*m.cols + cy
}

func (m Matrix) Get(x, y int) bool {
	return m.data[m.index(x, y)]
}

func (m *Matrix) Set(x, y int, val bool) {
	m.data[m.index(x, y)] = val
}
