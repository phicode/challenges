# Solution

```text
Part 1 input: aoc24/day15/example_small.txt
Result: 2028
Part 1: 40.386µs

Part 1 input: aoc24/day15/example.txt
Result: 10092
Part 1: 23.304µs

Part 1 input: aoc24/day15/input.txt
Result: 1426855
Part 1: 314.345µs

Part 2 input: aoc24/day15/example_p2.txt
Result: 618
Part 2: 14.327µs

Part 2 input: aoc24/day15/example.txt
Result: 9021
Part 2: 23.004µs

Part 2 input: aoc24/day15/input.txt
Result: 1404917
Part 2: 398.384µs
```

# Notes

## Part 2 Example
```text
             1111111111
   01234567890123456789
0  ####################
1  ##....[]....[]..[]##
2  ##............[]..##
3  ##..[][]....[]..[]##
4  ##....[]@.....[]..##
5  ##[]##....[]......##
6  ##[]....[]....[]..##
7  ##..[][]..[]..[][]##
8  ##........[]......##
9  ####################
```

Debug
```text
             1111111111
   01234567890123456789
0  ####################
1  ##....[]....[]..[]##
2  ##............[]..##
3  ##..[][]....[]..[]##
4  ##..[]........[]..##
5  ##[]##....[]......##
6  ##[]......[]..[]..##
7  ##..[][]..@.[][][]##
8  ##........[]......##
9  ####################
```