package gap

import "io"
import "testing"

func TestNewBuffer(t *testing.T) {

	var buf *Buffer = NewBuffer()
	if buf == nil {
		t.Errorf("A newly created buffer must not be nil!")
	}
}

func TestNewBufferLen(t *testing.T) {

	buf := NewBuffer()
	if buf.Len() != 0 {
		t.Errorf("A newly created buffer should have zero length!")
	}
}

func TestNewBufferCap(t *testing.T) {

	buf := NewBuffer()
	if buf.Cap() != 0 {
		t.Errorf("A newly created buffer should have zero capacity!")
	}
}

func TestNewBufferPos(t *testing.T) {

	buf := NewBuffer()
	if buf.Pos() != BufferStart {
		t.Errorf("The cursor of a newly created buffer should be at it's beginning!")
	}
}

const END_SIZE = 1 << 10

func TestResizedLen(t *testing.T) {

	buf := NewBuffer()
	buf.Resize(END_SIZE)

	if buf.Len() != 0 {
		t.Errorf("The buffer length should not change after resizing!")
	}
}

func TestResizedCap(t *testing.T) {

	buf := NewBuffer()
	buf.Resize(END_SIZE)

	if buf.Cap() != END_SIZE {
		t.Errorf("The buffer capacity should change after resizing!")
	}
}

var TEXT = []byte("Lorem ipsum dolor sit amet")

func TestWrite(t *testing.T) {

	buf := NewBuffer()
	n, err := buf.Write(TEXT)

	notAll, notNil := n != len(TEXT), err != nil

	if notAll {
		t.Errorf("A Write call should always write all the bytes given")

		if !notNil {
			t.Errorf("A Writer must return an error when not all the bytes were written")
		}
	}
	if notNil {
		t.Errorf("A Write call should never return a non-nil error")
	}
}

func TestPosAfterWrite(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	if buf.Pos() != len(TEXT) {
		t.Errorf("The cursor position should move by the number of bytes written")
	}
}

func TestLenAfterWrite(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	if buf.Len() != len(TEXT) {
		t.Errorf("The buffer length should be equal to the number of bytes written into it.")
	}
}

func TestCapAfterWrite(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	if buf.Cap() < len(TEXT) {
		t.Errorf("The buffer capacity should be no less than the number of bytes written to it.")
	}
}

func TestPosAfterTwoWrites(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.Write(TEXT)

	if buf.Pos() != len(TEXT)+len(TEXT) {
		t.Errorf("The cursor position should move by the number of bytes written")
	}
}

func TestLenAfterTwoWrites(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.Write(TEXT)

	if buf.Len() != len(TEXT)+len(TEXT) {
		t.Errorf("The buffer length should be equal to the number of bytes written into it.")
	}
}

func TestCapAfterTwoWrites(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.Write(TEXT)

	if buf.Cap() < len(TEXT)+len(TEXT) {
		t.Errorf("The buffer capacity should be no less than the number of bytes written to it.")
	}
}

func testPosAfterMoveToStart(t *testing.T, buf *Buffer) {

	buf.MoveTo(BufferStart)

	if buf.Pos() != BufferStart {
		t.Errorf("The cursors should be at the beginning of the buffer")
	}
}

func testPosAfterMoveBeforeStart(t *testing.T, buf *Buffer) {

	buf.MoveTo(BufferStart - 1)

	if buf.Pos() < BufferStart {
		t.Errorf("A cursor cannot move before the buffer's start")
	}
}

func testPosAfterMoveToEnd(t *testing.T, buf *Buffer) {

	buf.MoveTo(buf.Len())

	if buf.Pos() != buf.Len() {
		t.Errorf("The cursor should be at the end of the buffer")
	}
}

func testPosAfterMovePastEnd(t *testing.T, buf *Buffer) {

	buf.MoveTo(buf.Len() + 1)

	if buf.Pos() > buf.Len() {
		t.Errorf("The cursor should never move past the buffer's length")
	}
}

func TestMovementPastEmptyBufferEdges(t *testing.T) {

	buf := NewBuffer()

	testPosAfterMoveToStart(t, buf)
	testPosAfterMoveBeforeStart(t, buf)
	testPosAfterMoveToEnd(t, buf)
	testPosAfterMovePastEnd(t, buf)
}

func TestMovementPastNonEmptyBufferEdges(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	testPosAfterMoveToStart(t, buf)
	testPosAfterMoveBeforeStart(t, buf)
	testPosAfterMoveToEnd(t, buf)
	testPosAfterMovePastEnd(t, buf)
}

func TestMoveToPosHasZeroDelta(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	buf.MoveTo(buf.Len() / 2)
	if delta := buf.MoveTo(buf.Pos()); delta != 0 {
		t.Errorf("Moving to the current position should have zero displacement")
	}
}

func TestMoveAwayDelta(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	oldPos := buf.Pos()
	newPos := buf.Len() / 2

	delta := buf.MoveTo(buf.Len() / 2)
	if delta == 0 {
		t.Errorf("Moving to a new, valid position should have a non-zero delta")
	} else if delta != newPos-oldPos {
		t.Errorf("A move between valid positions should have the proper delta")
	}
}

func TestDontMoveLeftFromEmptyBufferStart(t *testing.T) {

	buf := NewBuffer()
	buf.MoveTo(BufferStart)
	delta := buf.MoveBy(-1)

	if delta != 0 {
		t.Errorf("A cursor cannot move left of the buffer's start")
	}
}

func testMoveToCurrentPosition(t *testing.T, buf *Buffer) {

	delta := buf.MoveTo(buf.Pos())

	if delta != 0 {
		t.Errorf("The delta should be zero when moving to the current position")
	}
}

func TestMoveToCurrentPosition(t *testing.T) {

	buf := NewBuffer()
	testMoveToCurrentPosition(t, buf)

	buf.Write(TEXT)
	testMoveToCurrentPosition(t, buf)

	buf.MoveTo(BufferStart)
	testMoveToCurrentPosition(t, buf)

	buf.MoveTo(buf.Len() / 2)
	testMoveToCurrentPosition(t, buf)
}

func TestDontMoveRightFromEmptyBufEnd(t *testing.T) {

	buf := NewBuffer()
	buf.MoveTo(buf.Len())
	delta := buf.MoveBy(1)
	if delta != 0 {
		t.Errorf("A cursor cannot move right of the buffer's end")
	}
}

func TestReadFromEmpty(t *testing.T) {

	buf, p := NewBuffer(), make([]byte, len(TEXT))

	n, err := buf.Read(p)
	if n != 0 {
		t.Error("Nothing should be read from an empty buffer")
	}
	if err != io.EOF {
		t.Error("A Read from an empty buffer should return io.EOF")
	}
}

func TestReadFromEmptyIntoEmpty(t *testing.T) {

	buf, p := NewBuffer(), []byte(nil)

	n, err := buf.Read(p)
	if n != 0 {
		t.Error("Nothing should be read from an empty buffer")
		t.Error("Nothing should be read into an empty slice")
	}
	if err != io.EOF {
		t.Error("A Read from an empty buffer should return io.EOF")
	}
}

func TestReadAll(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(BufferStart)

	p := make([]byte, buf.Len()*3/2)
	n, err := buf.Read(p)
	if n != buf.Len() {
		t.Error("The contents of the entire buffer should've been read")
	}
	if err != io.EOF {
		t.Error("Reading the entire buffer contets should return io.EOF")
	}
}

func TestReadAllAndFill(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(BufferStart)

	p := make([]byte, buf.Len())
	n, err := buf.Read(p)

	if n != len(p) {
		t.Error("The slice should be filled")
	}
	if err != io.EOF {
		t.Error("Reading the entire buffer contets should return io.EOF")
	}
}

func TestReadUntilEnd(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(buf.Len() / 2)
	startPos := buf.Pos()

	p := make([]byte, buf.Len())
	n, err := buf.Read(p)
	if n != buf.Len()-startPos {

		t.Error(
			"The number of bytes read should be equal the distance from the " +
				"starting position to the end",
		)
	}
	if err != io.EOF {
		t.Error("Reading up to the buffer's end should return io.EOF")
	}
}

func TestReadAndFill(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(buf.Len() / 3)

	p := make([]byte, buf.Len()/3)
	n, err := buf.Read(p)

	if n != len(p) {
		t.Error("The slice should be filled")
	}
	if err != nil {
		t.Error("There should be no I/O error")
	}
}

func TestReadUntilEndAndFill(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(buf.Len() / 3)

	p := make([]byte, buf.Len()-buf.Pos())
	n, err := buf.Read(p)

	if n != len(p) {
		t.Error("The slice should be filled")
	}
	if err != io.EOF {
		t.Error("Reading up to the buffer's end should return io.EOF")
	}
}

func TestReadIntoEmpty(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(buf.Len() / 2)

	p := []byte(nil)
	n, err := buf.Read(p)

	if n != 0 {
		t.Error("Nothing should be read into an empty slice")
	}
	if err != nil {
		t.Error("There should be no I/O error")
	}
}

func TestReadFromEnd(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	p := make([]byte, buf.Len())
	n, err := buf.Read(p)

	if n != 0 {
		t.Error("Nothing should be read from the end of a buffer")
	}
	if err != io.EOF {
		t.Error("Reading from a buffer's end should return io.EOF")
	}
}

func TestReadFromEndIntoEmpty(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)

	p := []byte(nil)
	n, err := buf.Read(p)

	if n != 0 {
		t.Error("Nothing should be read from the end of a buffer")
		t.Error("Nothing should be read into an empty slice")
	}
	if err != io.EOF {
		t.Error("Reading from a buffer's end should return io.EOF")
	}
}

func TestReadFromStartAndFill(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(BufferStart)

	p := make([]byte, buf.Len()/2)
	n, err := buf.Read(p)

	if n != len(p) {
		t.Error("The slice should be filled")
	}
	if err != nil {
		t.Error("There should be no I/O error")
	}
}

func TestReadFromStartIntoEmpty(t *testing.T) {

	buf := NewBuffer()
	buf.Write(TEXT)
	buf.MoveTo(BufferStart)

	p := []byte(nil)
	n, err := buf.Read(p)

	if n != 0 {
		t.Error("Nothing should be read into an empty slice")
	}
	if err != nil {
		t.Error("There should be no I/O error")
	}
}
