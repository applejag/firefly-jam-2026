package state

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"strconv"

	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"
)

type FieldError struct {
	Field string
	Err   error
}

func (err FieldError) Error() string {
	return err.Field + ": " + err.Err.Error()
}

func (err FieldError) Unwrap() error {
	return err.Err
}

type IndexedFieldError struct {
	Field string
	Index int
	Err   error
}

func (err IndexedFieldError) Error() string {
	return err.Field + "[" + strconv.Itoa(err.Index) + "]: " + err.Err.Error()
}

func (err IndexedFieldError) Unwrap() error {
	return err.Err
}

// Arbitrary number, mostly for future-proofing
const FileMarker byte = 0x37

type FieldState byte

const (
	FieldStateUnknown            FieldState = 0
	FieldStateFireflies          FieldState = 1
	FieldStateBattlesWonTotal    FieldState = 2
	FieldStateBattlesPlayedTotal FieldState = 3
	FieldStateMoney              FieldState = 4
)

var (
	_ encoding.BinaryAppender    = &GameState{}
	_ encoding.BinaryUnmarshaler = &GameState{}
)

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (g *GameState) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	if marker, err := buf.ReadByte(); err != nil || marker != FileMarker {
		return errors.New("invalid save file marker")
	}

	for {
		next, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch FieldState(next) {
		case FieldStateBattlesPlayedTotal:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "GameState.BattlesPlayedTotal", Err: err}
			}
			g.BattlesPlayedTotal = int(val)
		case FieldStateBattlesWonTotal:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "GameState.BattlesWonTotal", Err: err}
			}
			g.BattlesWonTotal = int(val)
		case FieldStateMoney:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "GameState.Money", Err: err}
			}
			g.Money = int(val)
		case FieldStateFireflies:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "GameState.Fireflies.len", Err: err}
			}
			if val > 99 {
				return FieldError{Field: "GameState.Fireflies.len", Err: errors.New("too many fireflies")}
			}
			g.Fireflies = make([]Firefly, val)
			for i := range int(val) {
				if err := g.Fireflies[i].UnmarshalBinaryBuf(buf); err != nil {
					return IndexedFieldError{Field: "GameState.Fireflies", Index: i, Err: err}
				}
			}
		default:
			return errors.New("unexpected game state field")
		}
	}
}

// AppendBinary implements [encoding.BinaryAppender].
func (g *GameState) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, FileMarker)

	b = append(b, byte(FieldStateFireflies))
	b = binary.AppendUvarint(b, uint64(len(g.Fireflies)))
	for i := range g.Fireflies {
		var err error
		b, err = g.Fireflies[i].AppendBinary(b)
		if err != nil {
			return b, IndexedFieldError{Field: "GameState.Fireflies", Index: i, Err: err}
		}
	}

	b = append(b, byte(FieldStateBattlesWonTotal))
	b = binary.AppendUvarint(b, uint64(g.BattlesWonTotal))
	b = append(b, byte(FieldStateBattlesPlayedTotal))
	b = binary.AppendUvarint(b, uint64(g.BattlesPlayedTotal))
	b = append(b, byte(FieldStateMoney))
	b = binary.AppendUvarint(b, uint64(g.Money))

	return b, nil
}

type FieldFirefly byte

const (
	FieldFireflyEOF           FieldFirefly = 0
	FieldFireflyID            FieldFirefly = 1
	FieldFireflyName          FieldFirefly = 2
	FieldFireflySpeed         FieldFirefly = 3
	FieldFireflyNimbleness    FieldFirefly = 4
	FieldFireflyBattlesPlayed FieldFirefly = 5
	FieldFireflyBattlesWon    FieldFirefly = 6
	FieldFireflyHat           FieldFirefly = 7
)

func (f *Firefly) UnmarshalBinaryBuf(buf *bytes.Buffer) error {
	for {
		next, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch FieldFirefly(next) {
		case FieldFireflyID:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.ID", Err: err}
			}
			f.ID = int(val)
		case FieldFireflyName:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.Name", Err: err}
			}
			f.Name = util.Name(val)
		case FieldFireflySpeed:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.Speed", Err: err}
			}
			f.Speed = int(val)
		case FieldFireflyNimbleness:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.Nimbleness", Err: err}
			}
			f.Nimbleness = int(val)
		case FieldFireflyBattlesPlayed:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.BattlesPlayed", Err: err}
			}
			f.BattlesPlayed = int(val)
		case FieldFireflyBattlesWon:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.BattlesWon", Err: err}
			}
			f.BattlesWon = int(val)
		case FieldFireflyHat:
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return FieldError{Field: "Firefly.Hat", Err: err}
			}
			f.Hat = int(val)
		case FieldFireflyEOF:
			return nil
		default:
			return errors.New("unexpected firefly field")
		}
	}
}

// AppendBinary implements [encoding.BinaryAppender].
func (f *Firefly) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, byte(FieldFireflyID))
	b = binary.AppendUvarint(b, uint64(f.ID))
	b = append(b, byte(FieldFireflyName))
	b = binary.AppendUvarint(b, uint64(f.Name))
	b = append(b, byte(FieldFireflySpeed))
	b = binary.AppendUvarint(b, uint64(f.Speed))
	b = append(b, byte(FieldFireflyNimbleness))
	b = binary.AppendUvarint(b, uint64(f.Nimbleness))
	b = append(b, byte(FieldFireflyBattlesPlayed))
	b = binary.AppendUvarint(b, uint64(f.BattlesPlayed))
	b = append(b, byte(FieldFireflyBattlesWon))
	b = binary.AppendUvarint(b, uint64(f.BattlesWon))
	b = append(b, byte(FieldFireflyHat))
	b = binary.AppendUvarint(b, uint64(f.Hat))

	return append(b, byte(FieldFireflyEOF)), nil
}
