## BlockID

- ブロック: ファイルの最小単位
- ブロックID: ファイル名とブロック番号を持つ
  - ファイル＝データベース,二次記憶装置
- ファイル名とブロック番号を持つ
  - あるファイルの中のデータの最小単位

```sh
ファイル: users.db
  ├── Block 0 → BlockID("users.db", 0)
  ├── Block 1 → BlockID("users.db", 1)
  ├── Block 2 → BlockID("users.db", 2)
```

## Page

- メモリでのデータの最小単位
- メモリ上のバッファとして扱われる
- このデータをファイル（ブロック）に書き込むことで永続化される
  - なので、ページサイズはブロックサイズと同じ
- int
  - 4byte
- byte
  - 先頭4byteがint値で、byteのデータ長を表す
  - なので、取得時は先頭4byteを読み取り、その値を元にbyteのデータ長を取得
  - set時は、byteのデータ長を先頭4byteに書き込み、その後ろにbyteのデータを書き込む
- string

## File

- ファイルをブロック単位で管理する
- ブロック: ファイルの最小単位
- ブロックを読んで Page に格納（Read）
  - 処理の流れ
    - blk.Number() * fm.blockSize で ブロックの開始位置（オフセット）を計算
    - Seek を使ってファイルの該当ブロックの位置に移動
    - その位置から Read で Page にデータを読み込む
📌 例（ブロックサイズ 4096 バイト）
- ブロック 0 の位置 = 0 * 4096 = 0 バイト
- ブロック 1 の位置 = 1 * 4096 = 4096 バイト
- ブロック 2 の位置 = 2 * 4096 = 8192 バイト
- Page をブロックに書き込む（Write）
- ファイルの末尾にブロックを追加（Append）
- ファイルサイズを計算し、ブロック数を計算（blockNum）

## LogMgr

- ログマネージャーはログレコードをログファイルに書き込むコンポーネント
  - ログの値はログレコードに保存され、ログレコードはログページで管理されログファイルに書き込まれる
  - ログレコードは後ろから追加される

- データ構造
- | boundary (4 bytes) | ....... records ........ |
- このboundary値は「ページ内の次にレコードを追加できる位置（オフセット）」を示す
- レコードはページの終端から先頭へ追加されるので boundary はレコード追加ごとに小さくなる
- イメージ図

## BufferPool

- pin
  - ブロックのデータをロックするために page(メモリ)に保持する

バッファプールの振る舞い

1. バッファプールは最初、すべてが unpin 状態
2. クライアントからリクエストが来るたびに、未使用のpageにブロックを読み込む
3. バッファプールのページがすべて埋まると置換が始まる
   1. ただし、pinされているページは置換できない
   2. そのため、unpinされるまで待たされる

- どの page を置換するかは、パフォーマンスに直結する
  - 理想的な置換戦略は、「今後一番長く使われないページを選ぶ」ことだが、「将来のアクセス」は分からない。
- クライアントはできるだけ早く unpin してあげるべき
