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
Na = Number of A presses
Nb = Number of B presses

X = Na*Xa + Nb*Xb
Y = Na*Ya + Nb*Yb
Cost = Na*3 + Nb*1 

# First Machine
X: 8400=Na*94+Nb*22
Y: 5400=Na*34+Nb*67

Maximise Nb
X: Na=8400/22 = 381.81.. = 381
Y: Na=5400/67 = 80.59.. = 80
min(381, 80) = 80

8400=Na*94+80*22
(8400-(80*22))/94=Na

6640=Na*94
Na=6640/94=70.63..
=> not an even integer => reduce Nb

8400=Na*94+79*22
8400=Na*94+1738
6640=Na*94
Na=6640/94=70.63..
```

```text
X = Na*Xa + Nb*Xb
Y = Na*Ya + Nb*Yb

Na*Xa + Nb*Xb - X = Na*Ya + Nb*Yb - Y
Na*Xa + Nb*Xb = Na*Ya + Nb*Yb - Y + X
Na*Xa - Na*Ya = Nb*Yb - Nb*Xb - Y + X
Na * (Xa - Ya) = Nb*Yb - Nb*Xb - Y + X
Na = (Nb*Yb - Nb*Xb - Y + X) / (Xa - Ya)


Na*94+Nb*22-8400 = Na*34+Nb*67-5400
Na*94+Nb*22 = Na*34+Nb*67-5400+8400

Na = (Nb*Yb - Nb*Xb - Y + X) / (Xa - Ya)
with
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400
Na = (Nb*67 - Nb*22 - 5400 + 8400) / (94 - 34)
Na = (Nb*45 + 3000) / 60
60 Na = (Nb*45 + 3000)

d = Nb-Na





```

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
