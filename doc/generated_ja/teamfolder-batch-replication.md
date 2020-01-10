# teamfolder batch replication 

Batch replication of team folders

# 利用方法

このドキュメントは"デスクトップ"フォルダを例として使用します.

## 実行

Windows:

```powershell
cd $HOME\Desktop
.\tbx.exe teamfolder batch replication -file TEAMFOLDER_NAME_LIST.csv
```

macOS, Linux:

```bash
$HOME/Desktop/tbx teamfolder batch replication -file TEAMFOLDER_NAME_LIST.csv
```

macOS Catalina 10.15以上の場合: macOSは開発者情報を検証します. 現在、`tbx`はそれに対応していません. 実行時の最初に表示されるダイアログではキャンセルします. 続いて、”システム環境設定"のセキュリティーとプライバシーから一般タブを選択します.
次のようなメッセージが表示されています:
> "tbx"は開発元を確認できないため、使用がブロックされました。

"このまま開く"というボタンがあります. リスクを確認の上、開いてください. ２回目の実行ではダイアログに"開く”ボタンがありますので、これを選択します

## オプション

| オプション       | 説明                                      | デフォルト |
|------------------|-------------------------------------------|------------|
| `-dst-peer-name` | Destination team account alias            |            |
| `-file`          | Data file for a list of team folder names |            |
| `-src-peer-name` | Source team account alias                 |            |

共通のオプション:

| オプション      | 説明                                                                                             | デフォルト     |
|-----------------|--------------------------------------------------------------------------------------------------|----------------|
| `-bandwidth-kb` | コンテンツをアップロードまたはダウンロードする際の帯域幅制限(Kバイト毎秒)0の場合、制限を行わない | 0              |
| `-concurrency`  | 指定した並列度で並列処理を行います                                                               | プロセッサー数 |
| `-debug`        | デバッグモードを有効にする                                                                       | false          |
| `-low-memory`   | 省メモリモード                                                                                   | false          |
| `-proxy`        | HTTP/HTTPS プロクシ (ホスト名:ポート番号)                                                        |                |
| `-quiet`        | エラー以外のメッセージを抑制し、出力をJSONLフォーマットに変更します                              | false          |
| `-secure`       | トークンをファイルに保存しません                                                                 | false          |
| `-workspace`    | ワークスペースへのパス                                                                           |                |

# ファイル書式

## 書式: File 

| 列   | 説明                | 値の説明 |
|------|---------------------|----------|
| name | Name of team folder | Sales    |

最初の行はヘッダ行です. プログラムはヘッダ行がない場合も認識します.

```csv
name
Sales
```

# ネットワークプロクシの設定

プログラムはシステム設定から自動的にプロクシ設定情報を取得します. しかしながら、それでもエラーが発生する場合には明示的にプロクシを指定することができます. `-proxy` オプションを利用します, `-proxy ホスト名:ポート番号`のように指定してください. なお、現在のところ認証が必要なプロクシには対応していません.

# 実行結果

作成されたレポートファイルのパスはコマンド実行時の最後に表示されます. もしコマンドライン出力を失ってしまった場合には次のパスを確認してください. [job-id]は実行の日時となります. このなかの最新のjob-idを各委任してください.

| OS      | Path                                                                                                      |
| ------- | --------------------------------------------------------------------------------------------------------- |
| Windows | `%HOMEPATH%\.toolbox\jobs\[job-id]\reports` (e.g. C:\Users\bob\.toolbox\jobs\20190909-115959.597\reports) |
| macOS   | `$HOME/.toolbox/jobs/[job-id]/reports` (e.g. /Users/bob/.toolbox/jobs/20190909-115959.597/reports)        |
| Linux   | `$HOME/.toolbox/jobs/[job-id]/reports` (e.g. /home/bob/.toolbox/jobs/20190909-115959.597/reports)         |

## レポート: verification 

レポートファイルは次の3種類のフォーマットで出力されます;
* `verification.csv`
* `verification.xlsx`
* `verification.json`

`-low-memory`オプションを指定した場合には、コマンドはJSONフォーマットのレポートのみを出力します.

レポートが大きなものとなる場合、`.xlsx`フォーマットのファイルは次のようにいくつかに分割されて出力されます;
`verification_0000.xlsx`, `verification_0001.xlsx`, `verification_0002.xlsx`...   

| 列         | 説明                                                                                                                                                                                           |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| diff_type  | 差分のタイプ`file_content_diff`: コンテンツハッシュの差分, `{left|right}_file_missing`: 左または右のファイルが見つからない, `{left|right}_folder_missing`: 左または右のフォルダが見つからない. |
| left_path  | 左のパス                                                                                                                                                                                       |
| left_kind  | フォルダまたはファイル                                                                                                                                                                         |
| left_size  | 左ファイルのサイズ                                                                                                                                                                             |
| left_hash  | 左ファイルのコンテンツハッシュ                                                                                                                                                                 |
| right_path | 右のパス                                                                                                                                                                                       |
| right_kind | フォルダまたはファイル                                                                                                                                                                         |
| right_size | 右ファイルのサイズ                                                                                                                                                                             |
| right_hash | 右ファイルのコンテンツハッシュ                                                                                                                                                                 |
