// Package mocks contém mocks compartilhados para testes unitários.
// Estes mocks são projetados para serem reutilizados em múltiplos arquivos de teste.
//
// Exemplo de uso:
//   storage := &mocks.MockStorage{}
//   storage.On("Upload", mock.Anything, mock.Anything).Return("path", nil)
//   svc := NewMyService(storage)
package mocks