package utils

type Grid[T any] struct {
	Slice  []T
	width  int
	height int
}

func GridFromDimensions[T any](width, height int) Grid[T] {
	slice := make([]T, width*height)
	return GridFromSlice(slice, width)
}

func GridFromSlice[T any](slice []T, width int) Grid[T] {
	return Grid[T]{slice, width, len(slice) / width}
}

func (grid *Grid[T]) Width() int {
	return grid.width
}

func (grid *Grid[T]) Height() int {
	return grid.height
}

func (grid *Grid[T]) GetCopy(position Point) T {
	return grid.Slice[position.ToIndex(grid.width)]
}

func (grid *Grid[T]) Get(position Point) *T {
	return &grid.Slice[position.ToIndex(grid.width)]
}

func (grid *Grid[T]) Set(position Point, value T) {
	grid.Slice[position.ToIndex(grid.width)] = value
}

func (grid *Grid[T]) Positions() GridPositionsIterator[T] {
	return GridPositionsIterator[T]{-1, grid}
}

func (grid *Grid[T]) Fill(value T) {
	for i := range grid.Slice {
		grid.Slice[i] = value
	}
}

func (grid *Grid[T]) Contains(point Point) bool {
	return point.X >= 0 && point.Y >= 0 && point.X < grid.width && point.Y < grid.height
}

func (grid *Grid[T]) PosFromIndex(index int) Point {
	return Point{index % grid.width, index / grid.width}
}

type GridPositionsIterator[T any] struct {
	index int
	grid  *Grid[T]
}

func (it *GridPositionsIterator[T]) Next() bool {
	it.index++
	return it.index < len(it.grid.Slice)
}

func (it *GridPositionsIterator[T]) Current() Point {
	return FromIndex(it.index, it.grid.width)
}
