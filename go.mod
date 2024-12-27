module github.com/amnn/adventofcode-2024

go 1.23.1

require (
	internal/grid v0.0.0
	internal/point v0.0.0
	internal/set v0.0.0
)

replace (
	internal/grid => ./internal/grid
	internal/point => ./internal/point
	internal/set => ./internal/set
)
