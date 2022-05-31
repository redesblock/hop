package soc

import (
	"github.com/redesblock/hop/core/swarm"
)

var _ swarm.ChunkValidator = (*Validator)(nil)

// SocVaildator validates that the address of a given chunk
// is a single-owner chunk.
type Validator struct {
}

// NewSocValidator creates a new SocValidator.
func NewValidator() swarm.ChunkValidator {
	return &Validator{}
}

// Validate performs the validation check.
func (v *Validator) Validate(ch swarm.Chunk) (valid bool) {
	s, err := FromChunk(ch)
	if err != nil {
		return false
	}

	address, err := s.Address()
	if err != nil {
		return false
	}
	return ch.Address().Equal(address)
}