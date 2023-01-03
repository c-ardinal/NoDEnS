# NoDEnS

NoDEnS is Diceroll Enhancement Service

Discord用のダイスボット

## 概要

DiscrodでBCDiceのコマンドを用いたダイスが振れるようになります。  
クトゥルフ神話TRPG(以下，CoC)用に作成したbotですが，
BCDiceがサポートしているシステムであれば利用可能です。

また，CoCのみで利用できる，Nodens独自の特殊なコマンドが有ります。

## 実行に必要なもの

- Go (v1.19.4)
- Ruby (v3.1)
- [ysakasin/bcdice-api](https://github.com/ysakasin/bcdice-api/tree/0.9.0) (v2.1.1)
- [DiscordBotのTokenとID](https://qiita.com/oosawa/items/e5b01e88a209d9087432)

※bcdice-apiはv1を使用します。v2への移行は近々行う予定です。

## 実行

1.Nodensのクローン

```sh
git clone https://github.com/c-ardinal/Nodens
cd Nodens
```

2.SystemConfig.jsonにDiscordBotのTokenとID，bcdice-apiのエンドポイントを記入

```json
{
    "discord-token": "Bot Hoge12Hoge34Foo56Bar78Baz90Boo",
    "discord-botid": "1234567890",
    "bcdice-endpoint": "http://localhost:9292/v1"
}
```

3.bcdice-apiを起動

4.Nodensをbuild&run

```sh
go get -d -v
go build
./Nodens
```

## BCDiceでサポートしているコマンド

[torgtaitai/BCDice](https://github.com/torgtaitai/BCDice)をご参照ください。

## Nodens独自にサポートしている共通コマンド

|コマンド      |引数                  |使用例                          |説明                                                                       |
|:------------:|:--------------------:|:------------------------------:|:-----------------------------------------------------------------------:|
|ShowVersion   |-                     |`ShowVersion`                   |バージョン情報を表示します                                                 |
|CreateSession |{SYSTEM_NAME}         |`CreateSession Cthulhu`         |左記のコマンドを実行したチャネルでセッションを生成し，ダイスボットを有効化します|
|CreateSession |--forced {SYSTEM_NAME}|`CreateSession --forced Cthulhu`|一度生成したセッションを破棄し，再生成します                                 |
|ConnectSession|{PARENT_CHANNEL_ID}   |`ConnectSession 1234567890`     |CreateSessionで生成したセッションに接続します                               |
|StoreSession  |-                     |`StoreSession`                  |CreateSessionで生成したセッションをサーバ上に保存します                      |
|Restore{SYSTEM_NAME}Session|-        |`RestoreCthulhuSession`         |StoreSessionで保存したセッションを復元します                                |

## Nodens独自にサポートしているCoC用コマンド

|コマンド      |引数                  |使用例                          |説明                                                             |演算機能の可否|演算機能の使用例|
|:------------:|:--------------------:|:------------------------------:|:---------------------------------------------------------------:|:------------:|:----------------:|
|regchara      |{CHARASHEET_URL}      |`regchara https://charasheet.vampire-blood.net/123456789abcdef`|キャラシートの情報を取得します    |✕            |-                 |
|check         |{ABILITY_SKILL_NAME}  |`check 聞き耳`                  |能力値もしくは技能値の初期値・開始値・現在値を確認出来ます       |✕            |-                 |
|ctrl          |{ABILITY_SKILL_NAME} {VAR_NUM} |`ctrl SAN -1`          |能力値もしくは技能値を加算/減算します                            |○            |`ctrl HP -(1d2+1)`|
|roll          |{ABILITY_SKILL_NAME}  |`roll 聞き耳`                   |S=5, F=95で1d100を振ります                                       |○            |`roll DEX*5`      |
|Sroll         |{ABILITY_SKILL_NAME}  |`Sroll 聞き耳`                  |S=5, F=95で1d100のシークレットダイスを振ります                   |○            |`Sroll 聞き耳+20` |
|sanc          |{SUCCESS} {FAIL}      |`sanc 1 1d2`                    |SAN値チェックロールおよびSAN値の加算/減算処理を実施します        |○            |`sanc 1 1d2+1`    |

## 各システム共通コマンドの使い方

### Case G-1. ダイスボットの有効化

```text
User: CreateSession Cthulhu
Bot: @User Session create successfully. (System: Cthulhu, ID: 1234567890)
```

### Case G-2. セッションの再生成

※同一チャネルでG-1が実行済み前提

```text
User: CreateSession Cthulhu
Bot: @User Session already exist.
User: CreateSession --forced Cthulhu
Bot: @User Session create successfully. (System: Cthulhu, ID: 1234567890)
```

### Case G-3. セッションの接続

※接続先チャネルでG-1もしくはG-2が実行済み前提  
※BotとのDMチャネルで使用することを想定しています

```text
User: ConnectSession 1234567890
Bot: @User Parent session connect successfully. (Parent: 1234567890, Child: 1357924680)
```

### Case G-4. BCDiceコマンドの実行

※同一チャネルでG-1もしくはG-2が実行済み前提

```text
User: 1d100
Bot: @User : (1D100) ＞ 39
User: 1d100<=50
Bot: @User : (1D100<=50) ＞ 71 ＞ 失敗
```

### Case G-5. シークレットダイスコマンドの実行

※子チャネルでG-3が実行済み前提  

```text
[Child Channel (BotとのDMチャネル)]
[4] User1: S1d100
[5] Bot: @User : (1D100) ＞ 36
```

```text
[Parent Channel (TRPGセッション用チャネル)]
[1] User1: hoge
[2] User2: fuga
[3] User3: piyo
[6] Bot: @User1 Secret dice.
```

シークレットダイスの内容開示機能は未実装です。  
内容を公開する場合は手動でお願いします。

## CoC専用コマンドの使い方

### Case C-1. キャラクターシート連携

※同一チャネルでG-1もしくはG-2が実行済み前提  
※キャラシートは [キャラクター保管所 様](https://charasheet.vampire-blood.net/) で作成し，保存後のURLを引数として指定してください

```text
User: regchara https://charasheet.vampire-blood.net/123456789abcdef
Bot: @User
====================
[名 前] ほげほげ ふーばー
[年 齢] 30歳
[性 別] 男
[職 業] 警察官
[ STR ] 9
[ CON ] 9
[ POW ] 12
[ DEX ] 13
[ APP ] 13
[ SIZ ] 16
[ INT ] 16 (Init: 15)
[ EDU ] 14
[ HP ] 13
[ MP ] 12
[ SAN ] 62 (Init: 60)
[ アイデア ] 80
[ 幸運 ] 60
[ 知識 ] 70
[メ モ] とある田舎の交番で働いている。
====================
```

### Case C-2. キャラクターシート連携ロール

※同一チャネルでC-1が実行済み前提

```text
User1: roll 知識
Bot: @User1 (1D100<=70) ＞ 61 ＞ 成功
User2: roll 知識
Bot: @User2 (1D100<=50) ＞ 97 ＞ 致命的失敗
```

### Case C-3. キャラクターシート連携シークレットロール

※子チャネルでG-3およびC-1が実行済み前提

```text
[Child Channel (BotとのDMチャネル)]
[3] User1: Sroll 回避
[4] Bot: @User1 : (1D100<=30) ＞ 61 ＞ 失敗
```

```text
[Parent Channel (TRPGセッション用チャネル)]
[1] User2: roll キック
[2] Bot: @User2 : (1D100<=60) ＞ 41 ＞ 成功
[5] Bot: @User1 Secret dice.
```

シークレットダイスの内容開示機能は未実装です。  
内容を公開する場合は手動でお願いします。

### Case C-4. 能力値の操作

※同一チャネルでC-1が実行済み前提

```text
User1: roll SAN
Bot: @User1 (1D100<=55) ＞ 61 ＞ 失敗
User1: 1d3
Bot: @User1 (1D3) ＞ 3
User1: ctrl SAN -3
Bot: @User1 [SAN] 55 => 52 (Diff: -3)

...

User1: roll SAN
Bot: @User1 (1D100<=52) ＞ 20 ＞ 成功
```

## その他

**Q. CoC以外のシステムに独自コマンドを実装する予定は有りますか?**
A. 無いです。対応させたいシステムがある場合は各自Forkなりプルリクをお願いします。
