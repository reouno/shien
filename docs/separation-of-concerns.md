# 疎結合設計の詳細

## 概要

Shienは、ログ記録機能とCLI機能が完全に独立して動作するよう設計されています。この文書では、各コンポーネントがどのように分離されているかを詳しく説明します。

## 独立したワークフロー

### 1. アクティビティログ記録のフロー

```
[Daemon Timer (5分ごと)]
    ↓
[Activity Worker (daemon.go:143)]
    ↓
[Activity Service] → [Activity Repository] → [SQLite]
```

- **完全に独立**: RPC層やCLI層を一切経由しない
- **自律的動作**: 5分ごとのタイマーで自動実行
- **エラー処理**: ログ記録の失敗は内部で処理され、他の機能に影響しない

### 2. CLIクライアントのフロー

```
[shienctl Command]
    ↓
[RPC Client] → [Unix Socket] → [RPC Server]
    ↓
[Service Layer] → [Repository] → [SQLite]
```

- **読み取り専用**: アクティビティログの記録には関与しない
- **独立した通信層**: Unix socketを使用した独自のプロトコル
- **サービス層経由**: ビジネスロジックはサービス層に集約

## レイヤー間の依存関係

### Service Layer (service/)
```go
type Services struct {
    Activity *ActivityService  // アクティビティ関連のビジネスロジック
    Config   *ConfigService    // 設定関連のビジネスロジック
}
```

**特徴**:
- Repository層のみに依存
- RPC層やWorker層から独立
- ビジネスロジックの単一の真実の源

### Repository Layer (database/repository/)
```go
type ActivityRepo struct {
    conn *sql.DB
}
```

**特徴**:
- データアクセスの抽象化
- SQLクエリの管理
- 上位層（Service/RPC/Worker）から独立

### RPC Layer (rpc/)
```go
type Server struct {
    services *service.Services  // サービス層への依存のみ
}
```

**特徴**:
- 通信プロトコルの処理に特化
- ビジネスロジックを持たない
- サービス層を通じてデータにアクセス

## 変更の影響範囲

### ログ記録ロジックの変更
| 変更内容 | 影響を受けるファイル | CLIへの影響 |
|---------|-------------------|------------|
| 記録間隔の変更 | daemon.go | なし |
| ログフォーマットの変更 | repository/activity.go | なし |
| 新しいログタイプの追加 | service/activity.go, repository/ | なし |

### CLIコマンドの変更
| 変更内容 | 影響を受けるファイル | ログ記録への影響 |
|---------|-------------------|---------------|
| 新コマンドの追加 | shienctl/main.go, rpc/protocol.go | なし |
| 出力フォーマットの変更 | shienctl/main.go | なし |
| 新しいクエリオプション | rpc/server.go, service/ | なし |

## 具体例

### 例1: ログ記録間隔を1分に変更
```go
// daemon/daemon.go の変更のみ
activityTicker := time.NewTicker(1 * time.Minute)  // 5分から1分に変更
```
→ CLI側は一切変更不要

### 例2: 新しいCLIコマンド「週次レポート」を追加
```go
// 1. protocol.go に新メソッド追加
const MethodGetWeeklyReport = "getWeeklyReport"

// 2. server.go にハンドラー追加
case MethodGetWeeklyReport:
    // 週次レポートのロジック

// 3. shienctl に新コマンド追加
```
→ ログ記録側は一切変更不要

## まとめ

この設計により、以下が実現されています：

1. **独立した開発**: ログ記録とCLI機能を別々のチームが並行開発可能
2. **影響の局所化**: 一方の変更が他方に波及しない
3. **テストの簡易化**: 各コンポーネントを独立してテスト可能
4. **将来の拡張性**: 新機能追加時の影響範囲が明確

この疎結合設計により、システムの保守性と拡張性が大幅に向上しています。