# arrtty
golangに似た文法の言語

# usage
```shell
# 基本
go run ./cmd/arrtty/main.go <filepath>
```
```shell
# フィボナッチ, n項目の値を終了コードとして返却
# デフォルトでn=10
go run ./cmd/arrtty/main.go ./examples/fib.txt
# exit code == 55
```

### 処理
- preprocess
  - tokenize : 文字列を分類し切り分ける
  - parse : 構文解析を行い読み込むことのできるコードか確認する(ここで生成されたnodeはassembleまで使用される)
  - analyze : 意味解析を行い型が一致しているかを確認する
- assemble
  - link : 意味解析された複数の意味ノード?を組み合わせ欠損のない意味ノードを作成する
  - compile : 欠損のない意味ノードからバーチャルマシン用の命令を作成する
- vm : 命令を実行するスタックマシン

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