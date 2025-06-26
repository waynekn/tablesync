package ws

import (
	"fmt"
	"strconv"
	"strings"
)

// mapToMatrix converts a map with keys in the format "row:col" to a 2D slice.
// It fills the matrix with values from the map, ensuring that each row has a length of `rowLen`.
// If a key's column index exceeds `rowLen`, that entry is skipped.
// If a key's format is invalid, an error is returned.
func mapToMatrix(input map[string]string, rowLen int) ([][]string, error) {
	var matrix [][]string
	lastRow := -1

	for key, val := range input {
		row, col, err := coordsFromString(key)
		if err != nil {
			return nil, err
		}
		if col >= rowLen {
			continue
		}
		if row > lastRow {
			extendMatrixRows(lastRow, row, rowLen, &matrix)
			lastRow = row
		}
		matrix[row][col] = val
	}

	return matrix, nil
}

// coordsFromString parses a string in the format "row:col" into two integers.
func coordsFromString(input string) (int, int, error) {
	coords := strings.Split(input, ":")
	if len(coords) != 2 {
		return -1, -1, fmt.Errorf("invalid input %q: expected format 'row:col'", input)
	}

	row, err := strconv.Atoi(coords[0])
	if err != nil {
		return -1, -1, fmt.Errorf("invalid row index in %q: %w", input, err)
	}

	col, err := strconv.Atoi(coords[1])
	if err != nil {
		return -1, -1, fmt.Errorf("invalid column index in %q: %w", input, err)
	}

	if row < 0 || col < 0 {
		return -1, -1, fmt.Errorf("negative indices in %q are not allowed", input)
	}

	return row, col, nil
}

// extendMatrixRows appends rows to the matrix until it reaches `desiredSize`
func extendMatrixRows(currSize int, desiredSize int, rowCap int, matrix *[][]string) {
	for i := 0; i < desiredSize-currSize; i++ {
		*matrix = append(*matrix, make([]string, rowCap))
	}
}
