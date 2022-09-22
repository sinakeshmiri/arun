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
	bin,err:=Make(src)
	if err!=nil{
		return "",err
	}
	return bin,nil
}
