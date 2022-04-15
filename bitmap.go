package xutils

import "strconv"

type BitMap struct {
	bits []uint64
	base int
}

// 最大数决定数组长度 65535/64 = 1024
// 存储数据大小 len(bits)*64

func NewBitMap(n int) *BitMap {
	return &BitMap{
		bits: make([]uint64, n),
		base: 0,
	}
}

func NewBitMapWithBase(n, b int) *BitMap {
	return &BitMap{
		bits: make([]uint64, n),
		base: b,
	}
}

var DefaultBitMap = NewBitMap(1024)

func (b *BitMap) index(n uint) (uint, uint) {
	return n / 64, n % 64
}

func (b *BitMap) grow(num uint) {
	var tbits []uint64
	if num+1 < uint(len(b.bits))+1024 {
		tbits = make([]uint64, len(b.bits)+1024)
	} else {
		tbits = make([]uint64, num+1)
	}
	copy(tbits, b.bits)
	b.bits = tbits
}

func (b *BitMap) Set(n int) {
	index, offset := b.index(uint(n - b.base))
	if index+1 > uint(len(b.bits)) {
		b.grow(index)
	}
	b.bits[index] ^= 1 << (offset) // 每一位代表一个数
}

func (b *BitMap) Include(n int) bool {
	index, offset := b.index(uint(n - b.base))
	if index+1 > uint(len(b.bits)) {
		return false
	}
	if b.bits[index]&(1<<offset) > 0 {
		return true
	}
	return false
}

func (b *BitMap) Del(n int) {
	index, offset := b.index(uint(n - b.base))
	if index+1 > uint(len(b.bits)) {
		return
	}
	b.bits[index] &^= 1 << offset
	if b.bits[len(b.bits)-1] != 0 {
		return
	}
	// 缩容
	i := len(b.bits) - 1
	for ; i >= 0; i-- {
		// 计算当前数组不为0的位置
		if b.bits[i] == 0 && i != len(b.bits)-1024 {
			continue
		}
		break
	}
	if i < len(b.bits)/2 || i == len(b.bits)-1024 {
		// 小于总组数一半或超过1023个,进行缩容
		b.bits = b.bits[:i+1]
	}
}

func (b *BitMap) Clear() {
	b.bits = make([]uint64, 0)
}

func (b *BitMap) All() []int {
	var rs []int
	for j := 0; j < len(b.bits); j++ {
		for i := 0; i < 64; i++ {
			if b.bits[j]&(1<<i) > 0 {
				rs = append(rs, j*64+i+b.base)
			}
		}
	}
	return rs
}

func (b *BitMap) String() string {
	var s string
	for j := 0; j < len(b.bits); j++ {
		for i := 0; i < 64; i++ {
			if b.bits[j]&(1<<i) > 0 {
				s += strconv.Itoa(j*64 + i + b.base)
				s += ","
			}
		}
	}
	if s != "" {
		return s[:len(s)-1]
	}
	return s
}
