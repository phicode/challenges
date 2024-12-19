# Solution

```text
Part 1 input: aoc24/day17/example.txt
Result: 4,6,3,5,6,3,5,2,1,0
Part 1: 44.885µs

Part 1 input: aoc24/day17/input.txt
Result: 1,3,5,1,7,2,5,1,6
Part 1: 14.557µs

Part 2 input: aoc24/day17/example2.txt
Result: 117440
Part 2: 1.483188ms
```

# Notes

- 3 bit computer
- 3 Register A, B, C of arbitrary precision

```text
Program: 2,4,1,3,7,5,4,7,0,3,1,5,5,5,3,0

2,4 BST: B = A % 8
1,3 BXL: B = B ^ 3
7,5 CDV: C = A / 2^B
4,7 BXC: B = B ^ C
0,3 ADV: A = A / 2^3
1,5 BXL: B = B ^ 5
5,5 OUT: Out = B % 8
3,0 if A!=0: IP=0  => A==0 terminates program

do {
  B = A % 8     // B=0-7
  B = B ^ 3     // flip bottom 2 bits; B=0-7
  C = A / 2^B   // C= A/(1,2,4,8,16,32,64,128)
  B = B ^ C
  A = A / 2^3  // A = A / 8
  B = B ^ 5
  Out = B % 8
} while(A != 0)

A=35184372088832
B = A % 8 = 0
C = A / 2^B = A / 1 = 35184372088832
B = B ^ C = 0^C = 35184372088832
A = A / 8 = 4398046511104
B = B ^ 5 = 35184372088837
Out = B % 8 = 5
```




```text
B^3
000 ^ 011 = 011 => 0 -> 3
001 ^ 011 = 010 => 1 -> 2
010 ^ 011 = 001 => 2 -> 1 
011 ^ 011 = 000 => 3 -> 0
100 ^ 011 = 111 => 4 -> 7
101 ^ 011 = 110 => 5 -> 6
110 ^ 011 = 101 => 6 -> 5
111 ^ 011 = 100 => 7 -> 4
```

```text
Part 2 input: aoc24/day17/example2.txt
input-a 0 registers: 0 0 0 1
input-a 64 registers: 1 0 0 2
input-a 192 registers: 0 0 0 3
input-a 2752 registers: 0 0 0 4
input-a 19136 registers: 0 0 0 5
input-a 117440 registers: 0 0 0 6
Result: 117440
Part 2: 1.709427ms
```


out len
a=0: 1
a=8: 2
a=64: 3
a=512: 4
a=4096: 5

8^1=8
8^2=64

3 correct @ 35520

4 correct @ 51904
117440-51904=65536
