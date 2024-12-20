package lib

type Field struct {
	n    int
	data []uint64
}

func NewField(n int) *Field {
	ints := n / 64
	if ints*64 < n {
		ints++
	}
	return &Field{
		n:    n,
		data: make([]uint64, ints),
	}
}

const (
	BucketMask  = 0xFFFF_FFFF_FFFF_FFC0
	BucketShift = 6 // 2^6=64
	ValueMask   = 0x3F
)

func (f *Field) Set(x int) {
	if x < 0 || x >= f.n {
		return
	}
	bucket := (uint64(x) & BucketMask) >> BucketShift
	bit := uint64(1 << (uint64(x) & ValueMask))
	f.data[bucket] |= bit
}

func (f *Field) FindMissingField() (int, bool) {
	for bucket, value := range f.data {
		if value != 0xFFFF_FFFF_FFFF_FFFF {
			x := bucket * 64
			for i := 0; i < 64; i++ {
				if value&(1<<i) == 0 {
					x += i
					if x >= f.n {
						return 0, false
					}
					return x, true
				}
			}
			panic("invalid state")
		}
	}
	return 0, false
}
