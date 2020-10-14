package core

import "context"

type NoOpTextEncrypt struct{}

func NewNoOpTextEncrypt() *NoOpTextEncrypt {
	return &NoOpTextEncrypt{}
}

func (m NoOpTextEncrypt) DecryptText(_ context.Context, cypherText string) (text string, err error) {
	return cypherText, nil
}

func (NoOpTextEncrypt) EncryptText(_ context.Context, text string) (cypherText string, err error) {
	return text, nil
}
