# arrtty
golangに互換 のある形を取りたい

### 処理(予定)
1. tokenize -> token chainを生成  
文字列を規則に従い切り分ける
2. parse -> syntax node treeを生成  
トークンが構文規則に従っているかをチェック
3. analyze -> semantic node treeを生成  
各ノードの型関係を調べながら扱いやすい形に整形する
4. assemble -> operation chainを生成  
vmに対する命令に変換する
5. link -> linked operation chainを生成  
所謂リンカのように複数のoperation chainを組み合わせて欠損を埋め、vmが実行可能な命令を作る
6. vm -> result  
最終的な演算を行い結果を返却する

### 文法
```text

program = toplevel*

toplevel = comment
         | "func" ident "(" funcParams? ")" funcReturns? stmt
         | "import" string
         | "var" ident types ("=" andor)?

stmt = expr
     | "return" expr?
     | "if" expr stmt ("else" stmt)?
     | "for" (expr? expr? expr?)? stmt
     | comment
     | "{" stmt* "}"

expr = assign

assign = "var" ident types ("=" andor)?
       | ident ":=" andor
       | andor ("=" andor)?

andor = equality ("&&" equality | "||" equality)*

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
        | float
        | string
        | bool
        | nil

types = "int" | "float" | "string" | "bool"
      | ident

callArgs = expr ("," expr)*

funcParams = ident types ("," ident types)*

funcReturns = types
            | "(" types ("," types)+ ")"

```