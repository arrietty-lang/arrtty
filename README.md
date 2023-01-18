# arrtty
golangに互換 のある形を取りたい

### 文法
```text

program = toplevel*

toplevel = comment
         | "func" ident "(" funcParams? ")" funcReturns? "{" stmt "}"
         | "import" string
         | "var" ident types ("=" andor)?

stmt = expr
     | "return" expr?
     | "if" expr "{" stmt "}" ("else" "{" stmt "}")?
     | "for" (expr? expr? expr?)? "{" stmt "}"
     | comment
     | stmt*

expr = assign

assign = "var" ident types ("=" andor)?
       | ident ":=" andor
       | andor ("=" andor)?

andor = equality ("==" equality | "!=" equality)*

equality = relational ("==" relational | "!=" relational)*

relational = add ("<" add | "<=" add + ">" add | ">=" add)*

add = mul ("+" mul | "-" mul)*

mul = unary ("*" unary | "/" unary | "%" unary)*

unary = ("+" | "-" | "!")? primary

primary = access

access = literal 

literal = "(" expr ")"
        | ident ("(" callArgs? ")")?
        | int
        | string
        | bool
        | null

types = "int" | "float" | "string" | "bool"
      | ident

callArgs = expr ("," expr)*

funcParams = ident types ("," ident types)*

funcReturns = types
            | "(" types ("," types)+ ")"

```