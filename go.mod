module okuzawats.com/go

go 1.19

require local.packages/repl v0.0.0

require (
	local.packages/lexer v0.0.0 // indirect
	local.packages/token v0.0.0 // indirect
)

replace local.packages/ast => ./ast

replace local.packages/lexer => ./lexer

replace local.packages/parser => ./parser

replace local.packages/repl => ./repl

replace local.packages/token => ./token
