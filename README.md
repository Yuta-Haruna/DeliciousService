# パン情報取得CLI(DeliciousService)

##  概要
contentful APIを利用し、FireStoreに下記情報を格納するCLIツール

・蜂蜜豆乳クランベリー

・黒ゴマポテロール

・黒七味と岩塩のフォカッチャ

##  使用技術
Windows 11
go version go1.19

フレームワーク：
 cobra

DB:
 firestore

##  CLIの使い方
1. このプロジェクトをクローンします。

```sh 
git clone https://github.com/Yuta-Haruna/DeliciousService.git
```

2. 認証情報のファイルを添付します。

   認証情報ファイルは、別途用意ください。(一部の方には配布しております。)
   
   認証ファイルを下記に配置ください。
 ```sh
DeliciousService
├─.gitignore
├─credentials.json　// ★ファイルを格納
└─main.go
```

3. プロジェクトのターミナルにて、下記コマンドを実行します。(アクセストークンは入れ替えてください。)
```sh
   go run main.go getData -t アクセストークン
```

4. 下記が出力されば、成功です。
```sh
 処理が完了しました
```
　
