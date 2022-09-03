This is regular expression compiler. This compiles regex to a deterministic finite automaton.

# Regular Expression Grammer
`symbol` and charactors enclosed by `'` are terminal.

```
Sum     ::= Concat '|' Sum
          | Concat

Concat  ::= Star Concat
          | Star

Star    ::= Primary '*'
          | Primary

Primary ::= Group
          | symbol

Group   ::= '(' Sum ')'
```
