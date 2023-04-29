# arrtty
golangに似た文法の言語(一部互換あり)

### 処理
- preprocess
  - tokenize : 文字列を分類し切り分ける
  - parse : 構文解析を行い読み込むことのできるコードか確認する(ここで生成されたnodeはassembleまで使用される)
  - analyze : 意味解析を行い型が一致しているかを確認する
- assemble
  - link : 意味解析された複数の意味ノード?を組み合わせ欠損のない意味ノードを作成する
  - compile : 欠損のない意味ノードからバーチャルマシン用の命令を作成する
- vm : 命令を実行するスタックマシン

### 機能
- vm
  - メイン関数から整数を終了コードとして返却できる
  - 関数を呼び出せる
  - 関数に引数を渡すせる
  - 関数からの戻り値を取得できる

### 文法
```text

program = toplevel*

toplevel = comment
         | "func" ident "(" funcParams? ")" funcReturns? stmt
         | "import" string
         | "var" ident types ("=" andor)?

stmt = expr
     | "return" expr? ("," expr)*
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

access = (ident ".")* literal 

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

### VM

- [x] NOP
---
- [x] ADD
- [x] SUB
---
- [x] CMP
- [x] LT
- [x] GT
- [x] LE
- [x] GE
---
- [x] JMP
- [x] JZ
- [ ] JNZ
- [ ] JE
- [ ] JNE
- [ ] JL
- [ ] JLE
- [ ] JG
- [ ] JGE
---
- [x] CALL
- [x] RET
---
- [x] PUSH
- [x] POP
---
- [x] MOV
- [ ] MSG
- [ ] LEN
---
- [ ] SYSCALL
  - WRITE
    - [ ] STDOUT
    - [ ] STDERR
    - [ ] FILE
  - READ
    - [ ] STDIN
    - [ ] FILE
- [x] EXIT