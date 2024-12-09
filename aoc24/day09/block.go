package main

import (
	"github.com/phicode/challenges/lib/assert"
	"math"
	"strconv"
	"strings"
)

type Blocks struct {
	Size       int
	Start, End *Block
}

type Block struct {
	StartIdx   int
	FileId     int // -1 = free
	Len        int
	Next, Prev *Block
}

func ParseBlocks(line string) *Blocks {
	size := CountAll(line)
	input := []rune(line)
	start := new(Block)
	end := start
	nextFileId := 0
	for i := 0; i < len(input); i++ {
		num := int(input[i] - '0')
		if i == 0 {
			start.StartIdx = i
			start.FileId = nextFileId
			nextFileId++
			start.Len = num
			continue
		}
		current := new(Block)
		current.Prev = end
		end.Next = current
		end = current
		if i%2 == 1 {
			current.FileId = -1
		} else {
			current.FileId = nextFileId
			nextFileId++
		}
		current.StartIdx = i
		current.Len = num
	}
	return &Blocks{Size: size, Start: start, End: end}
}

func (in *Blocks) String() string {
	sb := strings.Builder{}
	cur := in.Start
	for cur != nil {
		if cur.FileId == -1 {
			sb.WriteString(strings.Repeat(".", cur.Len))
		} else {
			idxStr := strconv.Itoa(cur.FileId)
			sb.WriteString(strings.Repeat(idxStr, cur.Len))
		}
		cur = cur.Next
	}
	return sb.String()
}

func (in *Blocks) Split(f *Block) (*Block, *Block) {
	p, q := f.Split()
	if in.Start == f {
		in.Start = p
	}
	if in.End == f {
		in.End = q
	}
	return p, q
}

func (in *Blocks) ToFileIdsArray() []int {
	ids := make([]int, in.Size)
	block := in.Start
	idx := 0
	for block != nil {
		idx = block.FillIds(ids, idx)
		block = block.Next
	}
	return ids
}

////////////////////////////////////////////////////////////

func (b *Block) PrevFile(maxFileId int) *Block {
	cur := b.Prev
	for cur != nil {
		if cur.FileId != -1 && cur.FileId <= maxFileId {
			return cur
		}
		cur = cur.Prev
	}
	return nil
}

func (b *Block) Split() (*Block, *Block) {
	p, q := new(Block), new(Block)
	prev, next := b.Prev, b.Next
	p.Next = q
	q.Prev = p
	if prev != nil {
		prev.Next = p
		p.Prev = prev
	}
	if next != nil {
		next.Prev = q
		q.Next = next
	}
	b.Prev, b.Next = nil, nil
	return p, q
}

func (b *Block) FillIds(ids []int, idx int) int {
	for i := 0; i < b.Len; i++ {
		ids[i+idx] = b.FileId
	}
	return idx + b.Len
}

func Rearrange(block *Blocks) {
	file := block.LastFile(math.MaxInt)
	//fmt.Println(block.String())
	for file != nil {
		id := file.FileId
		RearrangeFile(block, file)
		//fmt.Println(block.String())
		file = block.LastFile(id - 1)
	}
}

func RearrangeFile(block *Blocks, file *Block) {
	free := block.FindFreeFor(file)
	if free == nil {
		return
	}
	if free.Len > file.Len {
		p, q := block.Split(free)
		p.StartIdx = free.StartIdx
		q.StartIdx = free.StartIdx + file.Len
		p.Len = file.Len
		q.Len = free.Len - file.Len
		q.FileId = -1
		free = p
	}
	assert.True(free.Len == file.Len)
	free.FileId = file.FileId
	file.FileId = -1
}

func (in *Blocks) LastFile(maxId int) *Block {
	cur := in.End
	for cur != nil {
		if cur.FileId != -1 && cur.FileId <= maxId {
			return cur
		}
		cur = cur.Prev
	}
	return nil
}

func (in *Blocks) FindFreeFor(b *Block) *Block {
	cur := in.Start
	for cur != nil {
		if cur.FileId == -1 && cur.Len >= b.Len {
			if cur.StartIdx > b.StartIdx {
				return nil // do not move files to the right
			}
			return cur
		}
		cur = cur.Next
	}
	return nil

}
