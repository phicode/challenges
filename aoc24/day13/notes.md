# Solution

```text
Part 1 input: aoc24/day13/example.txt
Result: 480
Part 1: 47.99µs

Part 1 input: aoc24/day13/input.txt
Result: 31065
Part 1: 525.735µs

Part 2 input: aoc24/day13/example.txt
Result: 875318608908
Part 2: 22.482µs

Part 2 input: aoc24/day13/input.txt
Result: 93866170395343
Part 2: 474.097µs
```

# Notes

```text
# First Machine
X: 8400=Na*94+Nb*22
Y: 5400=Na*34+Nb*67

multiply X with Ay (34)
multiply Y with Ax (94)
X: 285600=a*3196+b*748
Y: 507600=a*3196+b*6298

subtract original X from Y
X: 285600=a*3196+b*748
Y: 222000=b*5550
b = 222000/5550
b=40

X: 285600=a*3196+40*748
   255680=a*3196
```
