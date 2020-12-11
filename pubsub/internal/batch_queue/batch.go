package batchqueue

import (
	"bytes"
	"encoding/binary"
)

type batch struct {
	fallbackBuffer [][]byte
	maxSize, curSize int
	header [4]byte
	buffer *bytes.Buffer
}

func (b *batch) Push(bts []byte) bool {
	length := len(bts)

	if length == 0 {
		return true
	}

	newSize := b.curSize + length + 4

	if newSize > b.maxSize {
		return false
	}

	binary.LittleEndian.PutUint32(b.header[:], uint32(length))
	b.fallbackBuffer = append(b.fallbackBuffer, bts)
	b.buffer.Write(b.header[:])
	b.buffer.Write(bts)
	b.curSize = newSize

	return true
}

func (b *batch) Reset() {
	b.curSize = 0
	b.fallbackBuffer = b.fallbackBuffer[:0]
	b.buffer.Reset()
}

func (b *batch) Bytes() []byte {
	return b.buffer.Bytes()
}

func (b *batch) Fallback() [][]byte {
	return b.fallbackBuffer
}

func newBatch(buffer *bytes.Buffer, maxSize int) batch {
	return batch{
		maxSize: maxSize,
		buffer:  buffer,
	}
}

// implementation for io.ReadFrom
// I don't want to delete this code

//type singleBatch struct {
//	pub []byte
//	header [4]byte
//	size, maxSize int
//	index int
//}
//
//func (b *singleBatch) Read(buf []byte) (int, error) {
//	bufOffset := 0
//	length := len(b.pub)
//
//	if b.index < 4 {
//		binary.LittleEndian.PutUint32(b.header[:], uint32(length))
//		copiedLength := copy(buf[bufOffset:], b.header[b.index:])
//		bufOffset += copiedLength
//		b.index += copiedLength
//	}
//
//	data := b.pub[b.index-4:]
//	copiedLen := copy(buf[bufOffset:], data)
//	bufOffset += copiedLen
//
//	if copiedLen < len(data) {
//		b.index += copiedLen
//	} else {
//		b.index = 0
//	}
//
//
//	// TODO: needs check
//	if b.index == 0 {
//		return bufOffset, io.EOF
//	}
//
//	return bufOffset, nil
//}
//
//func newSingleBatch(pub publication.Publication, maxSize int) *singleBatch {
//	return &singleBatch{
//		pub:     pub,
//		header:  [4]byte{},
//		size:    len(pub.Data) + 4,
//		maxSize: maxSize,
//	}
//}
//
//type batch struct {
//	pubs []publication.Publication
//	header [4]byte
//	pubOffset, index int
//	size, maxSize int
//}
//
//func (b *batch) Read(buf []byte) (int, error) {
//	bufOffset := 0
//	for ; b.pubOffset < len(b.pubs); b.pubOffset++ {
//		pub := b.pubs[b.pubOffset]
//		length := len(pub.Data)
//		if b.index < 4 {
//			binary.LittleEndian.PutUint32(b.header[:], uint32(length))
//			copiedLength := copy(buf[bufOffset:], b.header[b.index:])
//			fmt.Println(b.index, copiedLength, length, b.header)
//			bufOffset += copiedLength
//			b.index += copiedLength
//			if copiedLength < len(b.header[b.index-copiedLength:]) {
//				break
//			}
//		}
//
//		data := pub.Data[b.index-4:]
//		copiedLen := copy(buf[bufOffset:], data)
//		bufOffset += copiedLen
//		if copiedLen < len(data) {
//			b.index += copiedLen
//			break
//		} else {
//			b.index = 0
//		}
//	}
//
//	if b.index == 0 && b.pubOffset >= len(b.pubs) - 1 {
//		return bufOffset, io.EOF
//	}
//
//	return bufOffset, nil
//}
//
//func (b *batch) push(pub publication.Publication) bool {
//
//	newSize := b.size + len(pub.Data) + 4
//
//	if newSize > b.maxSize {
//		return true
//	}
//
//	b.pubs = append(b.pubs, pub)
//	b.size = newSize
//
//	return false
//}
//
//func newBatch(maxSize int) *batch {
//	return &batch{
//		pubs: []publication.Publication{},
//		header:  [4]byte{},
//		maxSize: maxSize,
//		size: 0,
//	}
//}
//
//func (b *batch) reset() {
//	b.pubs = b.pubs[:0]
//	b.size = 0
//	b.pubOffset = 0
//	b.index = 0
//}

