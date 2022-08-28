This is regular expression compiler. This compiler compiles regex to a deterministic finite automaton.

# Regular Expression Grammer
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
