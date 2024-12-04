# Solution

```text
Part 1 input: aoc24/day02/example.txt
Safe: 2
Part 1: 114.257µs

Part 1 input: aoc24/day02/input.txt
Safe: 202
Part 1: 246.897µs

Part 2 input: aoc24/day02/example.txt
Safe: 4
Part 2: 14.868µs

Part 2 input: aoc24/day02/input.txt
Safe: 271
Part 2: 351.256µs
```

# Notes


```text
7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
```

```text

1 2 7 8 9
diffs:
 1 5 1 1
  4 -4 0

7 6 4 2 1
 -1 -2 -2 -1

1 3 2 4 5
 2 -1 2 1
   1 1 -1
   
Removing 3:
  1 2 4 5
   1 2 1
 
 
```
