# Solution

```text
Part 1 input: aoc24/day12/example.txt
Result: 1930
Part 1: 41.148µs

Part 1 input: aoc24/day12/input.txt
Result: 1304764
Part 1: 732.016µs
```

# Notes

```text
AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
```

811292 is too high


```text
Edge Combinations for diagonal between a and X
  ab
  cX

- 2 of a, b, c are on
- all 3 are off
- only diagonal is on
  ..  XX  X.  .X  X.
  .X  .X  XX  XX  .X

the other 3 cases do not count as an edge
  .X  ..  XX
  .X  XX  XX 

a b c e
0 0 0 1
0 0 1 0
0 1 0 0
0 1 1 1
1 0 0 1
1 0 1 1
1 1 0 1
1 1 1 0

edge = ! (
             (!a && !b &&  c) ||
             (!a &&  b && !c) ||
             (a && b && c)
          ) 

edge = (!a && !b && !c) ||
       (!a &&  b &&  c) 
```
