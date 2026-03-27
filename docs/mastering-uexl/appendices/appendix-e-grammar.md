# Appendix E: Grammar Reference

This appendix describes the UExL grammar in EBNF notation. It is intended for tool builders, IDE extension authors, and developers who need to understand exactly what expressions are syntactically valid.

---

## Notation

```
::=   definition
|     alternation
[x]   optional (zero or one)
{x}   repetition (zero or more)
(x)   grouping
'x'   literal terminal
```

---

## Top-Level

```
Program ::= PipeExpression { '|' PipeType [Alias] ':' PipeExpression }
          | PipeExpression

PipeType ::= Identifier   (* named pipe: map, filter, reduce, etc. *)
           | ε            (* empty = passthrough pipe ':' *)

Alias ::= 'as' '$' Identifier
```

---

## Expressions (by precedence, lowest to highest)

```
PipeExpression ::= ConditionalExpression

ConditionalExpression ::= LogicalOrExpression
                        | LogicalOrExpression '?' ConditionalExpression ':' ConditionalExpression

LogicalOrExpression ::= LogicalAndExpression { '||' LogicalAndExpression }

LogicalAndExpression ::= BitwiseOrExpression { '&&' BitwiseOrExpression }

BitwiseOrExpression ::= BitwiseXorExpression { '|' BitwiseXorExpression }

BitwiseXorExpression ::= BitwiseAndExpression { '~' BitwiseAndExpression }

BitwiseAndExpression ::= EqualityExpression { '&' EqualityExpression }

EqualityExpression ::= ComparisonExpression { ('==' | '!=' | '<>') ComparisonExpression }

ComparisonExpression ::= NullishExpression { ('<' | '>' | '<=' | '>=') NullishExpression }

NullishExpression ::= ShiftExpression { '??' ShiftExpression }

ShiftExpression ::= AdditiveExpression { ('<<' | '>>') AdditiveExpression }

AdditiveExpression ::= MultiplicativeExpression { ('+' | '-') MultiplicativeExpression }

MultiplicativeExpression ::= PowerExpression { ('*' | '/' | '%') PowerExpression }

PowerExpression ::= UnaryExpression [ ('**' | '^') PowerExpression ]   (* right-assoc *)

UnaryExpression ::= ('-' | '!' | '~') UnaryExpression
                  | MemberExpression

MemberExpression ::= PrimaryExpression { Accessor }

Accessor ::= '.' Identifier
           | '?.' Identifier
           | '[' Expression ']'
           | '?.[' Expression ']'
           | '[' Slice ']'
           | '(' ArgumentList ')'   (* function call on result *)

Slice ::= [Expression] ':' [Expression] [':' Expression]
```

---

## Primary Expressions

```
PrimaryExpression ::= Literal
                    | Identifier
                    | ScopeVariable
                    | FunctionCall
                    | ArrayLiteral
                    | ObjectLiteral
                    | '(' Expression ')'

ScopeVariable ::= '$' Identifier    (* $item, $index, $acc, $last, $window, $chunk *)

FunctionCall ::= Identifier '(' ArgumentList ')'

ArgumentList ::= [ Expression { ',' Expression } ]

ArrayLiteral ::= '[' [ Expression { ',' Expression } ] ']'

ObjectLiteral ::= '{' [ ObjectEntry { ',' ObjectEntry } ] '}'
ObjectEntry   ::= (Identifier | StringLiteral | NumberLiteral) ':' Expression
```

---

## Literals

```
Literal ::= NumberLiteral
          | StringLiteral
          | BooleanLiteral
          | NullLiteral

NumberLiteral  ::= DecimalInteger
                 | DecimalFloat
                 | HexInteger          (* 0x[0-9a-fA-F]+ *)
                 | BinaryInteger       (* 0b[01]+ *)
                 | OctalInteger        (* 0o[0-7]+ *)
                 | ScientificNotation  (* e.g. 1e6, 2.5e-3 *)

StringLiteral  ::= "'" { AnyCharExceptSingleQuote } "'"
                 | '"' { AnyCharExceptDoubleQuote } '"'

BooleanLiteral ::= 'true' | 'false'

NullLiteral    ::= 'null'

Identifier     ::= Letter { Letter | Digit | '_' }
```

---

## Comments

UExL does not support comments. Expressions are single-line constructs.

---

## Whitespace

Whitespace (spaces, tabs) between tokens is ignored. Line breaks are not supported — all expressions must be on a single logical line.

---

## Pipe Restrictions

The following constructs are **not** allowed inside sub-expressions (parentheses or function arguments):
- Pipe operators (`|map:`, `|filter:`, etc.)
- Pipe aliases (`as $alias`)

```uexl
# VALID: pipe at top level
orders |filter: $item.total > 100

# INVALID: pipe inside parentheses
(orders |filter: $item.total > 100)  # parse error

# INVALID: alias inside parentheses
orders |map as $x: ($x.price)  # alias in sub-expression is rejected
```

---

## Language Keywords

These identifiers are reserved and cannot be used as variable names:

- `true`
- `false`
- `null`
- `as`

---

## Token Summary

| Token class | Examples |
|-------------|---------|
| Number | `42` `3.14` `0xFF` `0b1010` `0o17` `1e6` |
| String | `'hello'` `"world"` |
| Boolean | `true` `false` |
| Null | `null` |
| Identifier | `product` `basePrice` `x1` |
| ScopeVar | `$item` `$index` `$acc` `$last` `$window` `$chunk` |
| Operator | `+` `-` `*` `/` `%` `**` `^` `&` `\|` `~` `<<` `>>` `&&` `\|\|` `!` `??` `<` `>` `<=` `>=` `==` `!=` `<>` `?` `:` |
| Dot | `.` |
| QuestionDot | `?.` |
| LeftBracket | `[` |
| RightBracket | `]` |
| QuestionLeftBracket | `?.[` |
| LeftParen | `(` |
| RightParen | `)` |
| LeftBrace | `{` |
| RightBrace | `}` |
| Comma | `,` |
| Pipe | `\|` (pipe separator when followed by identifier and `:`) |
| As | `as` |
| Colon | `:` |
