package parser

import "github.com/google/uuid"

func ParseToUUID(input string) (output *uuid.UUID, err error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return
	}
	output = &id
	return
}
