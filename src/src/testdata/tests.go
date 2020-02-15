package testdata

//go:generate go run bldy.build/bldy/src/srcutils/cmd/srcencoder *.src
//go:generate go run bldy.build/bldy/src/srcutils/cmd/gen-corpus  *.src "../lexer/corpus"
//go:generate go run bldy.build/bldy/src/srcutils/cmd/gen-corpus  *.src "../parser/corpus"
