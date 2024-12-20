package rowcol

import "fmt"

func PrintByteGrid(g Grid[byte]) {
	for _, row := range g.Data {
		fmt.Println(string(row))
	}
	fmt.Println()
}
