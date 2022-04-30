package model

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

type Header struct {
	Target string

	Guess    *Word
	Proposal *Word

	PeerID string
}

type Word struct {
	Chars []*Char
}

type Char struct {
	Salt string
	Hash string
}

func NewHeader(guess string, gSalts []string, proposal, peerID, target string) (*Header, error) {
	var pSalt []string
	for i := 0; i < len(proposal); i++ {
		pSalt = append(pSalt, RandomString(30))
	}
	pw, err := GetChars(proposal, pSalt)
	if err != nil {
		return nil, err
	}

	gw, err := GetChars(guess, gSalts)
	if err != nil {
		return nil, err
	}

	return &Header{
		Target: target,
		PeerID: peerID,

		Guess: &Word{
			Chars: gw,
		},
		Proposal: &Word{
			Chars: pw,
		},
	}, nil
}

func Verify(guess string, challenge *Word) ([]bool, error) {
	result := make([]bool, len(challenge.Chars))

	if len(guess) != len(challenge.Chars) {
		return result, nil
	}

	var salts []string
	for _, ch := range challenge.Chars {
		salts = append(salts, ch.Salt)
	}

	gch, err := GetChars(guess, salts)
	if err != nil {
		return result, err
	}

	for i, ch := range challenge.Chars {
		result[i] = gch[i].Hash == ch.Hash
	}

	return result, nil
}

var ErrSaltsAndCharsDidntMatch = errors.New("number of salts and number of letters didn't match")

func GetChars(word string, salts []string) ([]*Char, error) {
	if len(word) != len(salts) {
		return nil, ErrSaltsAndCharsDidntMatch
	}
	var chars []*Char
	h := sha256.New()
	for i, r := range word {
		salt := salts[i]
		h.Reset()
		h.Write([]byte{byte(r)})
		h.Write([]byte(salt))

		hash := fmt.Sprintf("%x", h.Sum(nil))[:45]

		ch := &Char{
			Salt: salt,
			Hash: hash,
		}

		chars = append(chars, ch)
	}

	return chars, nil
}