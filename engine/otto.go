package engine

import (
	"errors"

	"github.com/robertkrimen/otto"
)

type ottoEngine struct {
	otto *otto.Otto
}

func NewOttoEngine() Engine {
	return &ottoEngine{
		otto: otto.New(),
	}
}

func (e *ottoEngine) Clone() Engine {
	return &ottoEngine{
		otto: e.otto.Copy(),
	}
}

func (e *ottoEngine) Load(src []byte) error {
	_, err := e.otto.Run(src)
	return err
}

func (e *ottoEngine) RunReact(src string) (string, error) {
	v, err := e.otto.Run(src)
	if err != nil {
		return "", err
	}

	if !v.IsString() {
		return "", errors.New("Expected string result, actual type is: " + v.Class())
	}

	return v.String(), nil
}
