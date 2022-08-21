# golex: Toy implementation of Lexer Generator


## Regular Expression Grammer
```
Sum     ::= Concat '|' Sum
          | Concat

Concat  ::= Star Concat
          | Star

Star    ::= Primary '*'
          | Primary

Primary ::= Group
          | Symbol

Group   ::= '(' Sum ')'
```
