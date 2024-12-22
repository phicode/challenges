# Solution

```text

```

# Notes

```text
789
456
123
 0A
 
max move distance: x=2, y=3
7->A: 
  - >>vvv
  - >v>vv
  - >vv>v
  - >vvv>
  - v>>vv
  - v>v>v
  - v>vv>
  - vv>>v
  - vv>v>
  - vvv>>
10 combinations, of which 1 is invalid: vvv>> => 9 Combinations

want:
  - >>vvv x
  - >v>vv x
  - >vv>v x
  - >vvv> x
  - v>>>v -> missing
  - v>v>v x
  - v>vv> x
  - v>>vv x
  - vv>>v x
  - vv>v> x
```

cost Explosions:

move up:
  - ^A
  - <A>A

move right:
  - vA
  - vA^A

move down:
  - v<A
  - v<A<A>>^A

  - <vA
  - v<<A>A>^A

move left:
  - <v<A
  - <v<A>A<A>>^A

  - v<<A
  - v<A<AA>>^A

- 'up' and 'right' moves are equaly expensive on the direction pad
  - 2 moves '<A' and 'vA' respectively
- 'down' costs 3 moves 
