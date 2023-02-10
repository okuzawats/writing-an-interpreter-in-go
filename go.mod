module okuzawats.com/go

go 1.19

require (
    local.packages/lexer v0.0.0
    local.packages/repl v0.0.0
    local.packages/token v0.0.0
)

replace local.packages/lexer => ./lexer
replace local.packages/repl => ./repl
replace local.packages/token => ./token
