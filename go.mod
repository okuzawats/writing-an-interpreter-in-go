module okuzawats.com/go

go 1.23.2

require local.packages/ast v0.0.0 // indirect

require local.packages/evaluator v0.0.0 // indirect

require local.packages/object v0.0.0 // indirect

require local.packages/lexer v0.0.0 // indirect

require local.packages/parser v0.0.0 // indirect

require local.packages/repl v0.0.0

require local.packages/token v0.0.0 // indirect

replace local.packages/ast => ./ast

replace local.packages/evaluator => ./evaluator

replace local.packages/lexer => ./lexer

replace local.packages/object => ./object

replace local.packages/parser => ./parser

replace local.packages/repl => ./repl

replace local.packages/token => ./token
