# UExL Go

Universal Expression Language parser for go projects!

Operators supported:

- Logical Operators             && ||
- Bitwise Operators             & | ^
- Equality Operators            == !=
- Comparison Operators          <= >= < >
- Bitwise Shift Operators       << >>
- Additive Operators            + -
- Multiplicative Operators      * /
- Modulus Operator              %
- Dot Operator                  .
- Grouping Operator             ()

Operator Precedence:

Operators | Type | Associativity
--- | --- | ---
`(` `)` | Parentheses | Left to Right
`.` | Dot | Left to Right
`%` | Modulus | Left to Right
`*` `/` | Multiplicative | Left to Right
`+` `-` | Additive | Left to Right
`<<` `>>` | Bitwise Shift | Left to Right
`<` `>` `<=` `>=` | Comparison | Left to Right
`==` `!=` | Equality | Left to Right
`&` `\|` `^` | Bitwise | Left to Right
`&&` `\|\|` | Logical | Left to Right
`\|:` | Pipe | Left to Right
