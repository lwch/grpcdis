package obj

import (
	"bytes"
)

type dictBucket struct {
	tophash [8]uint8
	keys    [8][]byte
	values  [8]interface{}
	sets    uint8 // means each index have value
	next    *dictBucket
}

func (bucket *dictBucket) setIdx(i uint8, top uint8, key []byte, value interface{}) {
	bucket.tophash[i] = top
	bucket.keys[i] = key
	bucket.values[i] = value
	bucket.sets |= 1 << i
}

func (bucket *dictBucket) set(top uint8, key []byte, value interface{}) (new, ok bool) {
	idx := -1
	col := uint8(0)
	empty := 8
	for i := uint8(0); i < 8; i++ {
		if (bucket.sets>>i)&1 == 0 {
			col = i
		} else {
			empty--
		}
		if bucket.tophash[i] != top {
			continue
		}
		if bytes.Equal(bucket.keys[i], key) {
			idx = int(i)
			break
		}
	}
	if idx == -1 {
		if empty == 0 { // no column
			return false, false
		} else {
			bucket.setIdx(col, top, key, value)
			return true, true
		}
	}
	bucket.setIdx(uint8(idx), top, key, value)
	return false, true
}

func (bucket *dictBucket) index(top uint8, key []byte) int {
	for i := 0; i < 8; i++ {
		if bucket.tophash[i] != top {
			continue
		}
		if !bytes.Equal(bucket.keys[i], key) {
			continue
		}
		return i
	}
	return -1
}

func (bucket *dictBucket) add(top uint8, key []byte, value interface{}) bool {
	idx := -1
	col := uint8(0)
	empty := 8
	for i := uint8(0); i < 8; i++ {
		if (bucket.sets>>i)&1 == 0 {
			col = i
		} else {
			empty--
		}
		if bucket.tophash[i] != top {
			continue
		}
		if bytes.Equal(bucket.keys[i], key) {
			idx = int(i)
			break
		}
	}
	if idx == -1 {
		if empty == 0 {
			return false
		} else {
			bucket.setIdx(col, top, key, value)
			return true
		}
	}
	return false
}

func (bucket *dictBucket) search(top uint8, key []byte) (*dictBucket, int) {
	idx := bucket.index(top, key)
	if idx != -1 {
		return bucket, idx
	}
	if bucket.next != nil {
		return bucket.next.search(top, key)
	}
	return nil, -1
}

func (bucket *dictBucket) end() *dictBucket {
	if bucket.next == nil {
		return bucket
	}
	return bucket.next.end()
}

func (bucket *dictBucket) lead() bool {
	if bucket.next == nil {
		return false
	}
	bucket.tophash = bucket.next.tophash
	bucket.keys = bucket.next.keys
	bucket.values = bucket.next.values
	bucket.sets = bucket.next.sets
	bucket.next = bucket.next.next
	return true
}

// DictValue hash map value
type Dict struct {
	count     int
	b         uint8
	noverflow uint

	buckets    []dictBucket
	oldBuckets []dictBucket
	rehash     int

	nextOverflow *dictBucket
}

// NewDict create dict object
func NewDict() *Dict {
	const initSize = 1
	return &Dict{
		b:       initSize,
		buckets: make([]dictBucket, bucketShift(initSize)),
		rehash:  -1,
	}
}

func bucketShift(b uint8) uint64 {
	return 1 << b
}

func bucketMask(b uint8) uint64 {
	return bucketShift(b) - 1
}

func tophash(hash uint64) uint8 {
	return uint8(hash >> (8 * 7)) // high 8bit
}

// Bernstein hash
func hash(str []byte) uint64 {
	hash := uint64(5381)
	for i := 0; i < len(str); i++ {
		hash = (hash << 5) + hash + uint64(str[i])
	}
	return hash
}

// Set set value
func (dict *Dict) Set(key []byte, value interface{}) {
	hash := hash(key)
	if dict.growing() {
		dict.grow()
	}
	bucket := &dict.buckets[hash&bucketMask(dict.b)]
	top := tophash(hash)
	for {
		for i := 0; i < 8; i++ {
			if new, ok := bucket.set(top, key, value); ok {
				if new {
					dict.count++
				}
				return
			}
		}
		next := bucket.next
		if next == nil {
			break
		}
		bucket = next
	}
	if dict.nextOverflow != nil {
		dict.nextOverflow.setIdx(0, top, key, value)
		bucket.next = dict.nextOverflow
		dict.nextOverflow = dict.nextOverflow.next
	} else {
		newBucket := new(dictBucket)
		newBucket.setIdx(0, top, key, value)
		bucket.next = newBucket
	}
	dict.count++
	dict.noverflow++
	if dict.wantGrow() {
		dict.markGrow()
	}
}

// Get get value
func (dict *Dict) Get(key []byte) interface{} {
	if dict.count == 0 {
		return nil
	}
	if dict.growing() {
		dict.grow()
	}
	hash := hash(key)
	mask := bucketMask(dict.b)
	bucket := &dict.buckets[hash&mask]
	if dict.oldBuckets != nil {
		bucket = &dict.oldBuckets[hash&(mask>>1)]
	}
	top := tophash(hash)
	for ; bucket != nil; bucket = bucket.next {
		for i := 0; i < 8; i++ {
			if bucket.tophash[i] != top {
				continue
			}
			if bytes.Equal(bucket.keys[i], key) {
				continue
			}
			return bucket.values[i]
		}
	}
	return nil
}

// Del delete value
func (dict *Dict) Del(key []byte) {
}

func (dict *Dict) growing() bool {
	return dict.rehash != -1
}

func (dict *Dict) grow() {
	for {
		move := false
		oldBucket := &dict.oldBuckets[dict.rehash]
		for i := 0; i < 8; i++ {
			if (oldBucket.sets>>i)&1 == 1 {
				hash := hash(oldBucket.keys[i])
				top := tophash(hash)
				newBucket := &dict.buckets[hash&bucketMask(dict.b)]
				k := oldBucket.keys[i]
				v := oldBucket.values[i]
				if bk, _ := newBucket.search(top, k); bk == nil {
					newBucket = newBucket.end()
					if !newBucket.add(top, k, v) {
						if dict.nextOverflow != nil {
							dict.nextOverflow.setIdx(0, top, k, v)
							newBucket.next = dict.nextOverflow
							dict.nextOverflow = dict.nextOverflow.next
						} else {
							newBucket = new(dictBucket)
							newBucket.setIdx(0, top, k, v)
						}
						dict.noverflow++
					}
				}
				oldBucket.tophash[i] = 0
				oldBucket.keys[i] = nil
				oldBucket.values[i] = nil
				oldBucket.sets ^= (1 << i)
				move = true
				break
			}
		}
		if move {
			return
		} else if !oldBucket.lead() {
			dict.rehash++
			if dict.rehash >= len(dict.oldBuckets) {
				dict.rehash = -1
				dict.oldBuckets = nil
				return
			}
		}
	}
}

func (dict *Dict) wantGrow() bool {
	b := dict.b
	if b > 15 {
		b = 15
	}
	return dict.noverflow > (1 << (uint(b) & 15))
}

func (dict *Dict) markGrow() {
	dict.oldBuckets = dict.buckets
	dict.b++
	dict.buckets = make([]dictBucket, bucketShift(dict.b))
	dict.rehash = 0
}
