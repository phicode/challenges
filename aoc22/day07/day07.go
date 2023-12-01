package main

// https://adventofcode.com/2022/day/7

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	Process("aoc22/day07/example.txt")
	Process("aoc22/day07/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	commands := SplitByCommand(lines)

	fs := NewFileSystem()
	for i, cmd := range commands {
		fmt.Println("Command:", cmd[0])
		if i == 0 {
			if len(cmd) != 1 || cmd[0] != "$ cd /" {
				panic("unexpected first command")
			}
			continue
		}
		for i := 1; i < len(cmd); i++ {
			fmt.Println("> ", cmd[i])
		}
		fs.apply(cmd)
	}
	fs.print()

	sf := SizeFinder{Size: 100_000}
	fs.visit(sf.visit)
	fmt.Println("sum:", sf.Sum)

	free := 70000000 - fs.Root.Size()
	todelete := 30000000 - free
	fmt.Println("need additional free space:", todelete)
	deleteFinder := DeleteFinder{Delete: todelete}
	fs.visit(deleteFinder.visit)
	fmt.Println("delete dir with size:", deleteFinder.Node.Size())

	fmt.Println()
}

func ReadInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

func SplitByCommand(lines []string) [][]string {
	var commands [][]string
	var current []string
	for _, line := range lines {
		// start of new command
		if len(line) > 0 && line[0] == '$' {
			if len(current) > 0 {
				commands = append(commands, current)
				current = nil
			}
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		commands = append(commands, current)
	}
	return commands
}

////////////////////////////////////////////////////////////

type Node struct {
	Parent *Node
	Name   string
	Dirs   map[string]*Node
	Files  map[string]int
}

func (n *Node) print(indent int) {
	prefix := strings.Repeat(" ", indent*2)
	fmt.Printf("%s- %s (dir, size=%d)\n", prefix, n.Name, n.Size())
	for _, dir := range n.Dirs {
		dir.print(indent + 1)
	}
	for file, size := range n.Files {
		fmt.Printf("%s  - %s (file, size=%d)\n", prefix, file, size)
	}
}

func NewNode(parent *Node, name string) *Node {
	return &Node{
		Parent: parent,
		Name:   name,
		Dirs:   make(map[string]*Node),
		Files:  make(map[string]int),
	}
}

func (n *Node) Size() int {
	sum := 0
	for _, dir := range n.Dirs {
		sum += dir.Size()
	}
	for _, size := range n.Files {
		sum += size
	}
	return sum
}

func (n *Node) visit(v visitor) {
	for _, dir := range n.Dirs {
		dir.visit(v)
	}
	v(n)
}

////////////////////////////////////////////////////////////

type FileSystem struct {
	Root    *Node
	Current *Node
}

func NewFileSystem() *FileSystem {
	root := NewNode(nil, "/")
	return &FileSystem{
		Root:    root,
		Current: root,
	}
}

func (s *FileSystem) apply(cmd []string) {
	if cmd[0] == "$ ls" {
		s.applyLs(cmd[1:])
		return
	}

	if cmd[0] == "$ cd .." {
		if s.Current == nil || s.Current.Parent == nil {
			panic("invalid filesystem traversal")
		}
		if len(cmd) > 1 {
			panic("invalid cd .. length")
		}
		s.Current = s.Current.Parent
		return
	}

	if strings.HasPrefix(cmd[0], "$ cd ") {
		if len(cmd) > 1 {
			panic("invalid cd directory length")
		}
		dir := cmd[0][5:]
		if node, ok := s.Current.Dirs[dir]; ok {
			s.Current = node
		} else {
			fmt.Println("entering yet unknown directory", dir)
			node := NewNode(s.Current, dir)
			s.Current.Dirs[dir] = node
			s.Current = node
		}
		return
	}

	panic(fmt.Errorf("unknown command: %q", cmd[0]))
}

func (s *FileSystem) applyLs(files []string) {
	for _, file := range files {
		if strings.HasPrefix(file, "dir ") {
			dir := file[4:]
			s.Current.Dirs[dir] = NewNode(s.Current, dir)
			continue
		}
		var size int
		var name string
		n, err := fmt.Sscanf(file, "%d %s", &size, &name)
		if n != 2 || err != nil {
			panic(fmt.Errorf("invalid ls line: %q", file))
		}
		s.Current.Files[name] = size
	}
}

func (s *FileSystem) print() {
	s.Root.print(0)
}

type visitor func(n *Node)

func (s *FileSystem) visit(v visitor) {
	s.Root.visit(v)
}

////////////////////////////////////////////////////////////

type SizeFinder struct {
	Size  int
	Sum   int
	Nodes []*Node
}

func (sf *SizeFinder) visit(n *Node) {
	s := n.Size()
	if s <= sf.Size {
		sf.Sum += s
		sf.Nodes = append(sf.Nodes, n)
	}
}

////////////////////////////////////////////////////////////

type DeleteFinder struct {
	Delete int
	Node   *Node
}

func (df *DeleteFinder) visit(n *Node) {
	s := n.Size()
	if s < df.Delete {
		return
	}
	if df.Node == nil {
		df.Node = n
		return
	}
	if df.Node.Size() > s {
		df.Node = n
	}
}
