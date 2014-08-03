package gap

import "io"

const (
	BufferStart = 0
)

// A text Buffer
type Buffer struct {
	pos, len, cap int
	data          []byte
}

// Create a new Buffer
func NewBuffer() *Buffer {
	return &Buffer{}
}

// Cursor position within the Buffer
func (buf *Buffer) Pos() int {
	return buf.pos
}

// Length of the text within the Buffer (in bytes)
func (buf *Buffer) Len() int {
	return buf.len
}

// Capacity of the Buffer (in bytes)
func (buf *Buffer) Cap() int {
	return buf.cap
}

// Resize a Buffer to the given size
func (buf *Buffer) Resize(newCap int) {
	buf.cap = newCap
}

// Write a byte slice into the Buffer
func (buf *Buffer) Write(p []byte) (int, error) {
	buf.pos += len(p)
	buf.len += len(p)
	buf.cap += len(p)
	buf.data = append(buf.data, p...)
	return len(p), nil
}

// Move the cursor to the specified position. Doesn't allow the cursor to leave
// the Buffer. Returns the displacement.
func (buf *Buffer) MoveTo(newPos int) int {

	delta := 0
	if BufferStart <= newPos && newPos <= buf.Len() {
		delta = newPos - buf.pos
		buf.pos = newPos
	}
	return delta
}

// Move the cursor by the specified displacement. Doesn't allow the cursor to
// leave the Buffer. Returns the actual displacement.
func (buf *Buffer) MoveBy(delta int) int {
	return buf.MoveTo(buf.Pos() + delta)
}

func (buf *Buffer) Read(p []byte) (n int, err error) {

	n = copy(p, buf.data[buf.pos:])
	buf.pos += n

	if buf.Pos() == buf.Len() {
		err = io.EOF
	}

	return
}
