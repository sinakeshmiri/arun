package compiler

import (

)

// Adapter implements the DbPort interface
type Adapter struct {
}

// NewAdapter creates a new Adapter
func NewAdapter() (*Adapter, error) {
	return &Adapter{}, nil
}

// CloseDbConnection closes the db  connection
func (da Adapter) Compile(src string) (string, error) {
	return(Make(src))
}
