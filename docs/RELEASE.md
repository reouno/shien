# リリース手順

## クイックリリース（推奨）

```bash
# scripts/release.sh を使用
./scripts/release.sh 0.1.2
```

このスクリプトは以下を自動的に実行します：
1. タグの作成とプッシュ
2. GitHub Actionsの完了待機
3. SHA256の計算
4. Homebrew formulaの更新
5. formulaのコミット・プッシュ

## 手動リリース手順

### 1. タグを作成してプッシュ

```bash
git tag -a v0.1.2 -m "Release version 0.1.2"
git push origin v0.1.2
```

### 2. GitHub Actionsの完了を待つ

https://github.com/reouno/shien/actions でビルドの完了を確認

### 3. リリースファイルのSHA256を取得

```bash
curl -L https://github.com/reouno/shien/releases/download/v0.1.2/shien-darwin-arm64.tar.gz -o /tmp/shien.tar.gz
shasum -a 256 /tmp/shien.tar.gz
```

### 4. Homebrew formulaを更新

`homebrew-shien`リポジトリで：

```bash
cd ~/homebrew-shien  # または適切なパス
```

`Formula/shien.rb`を編集：
- `version "0.1.1"` → `version "0.1.2"`
- URLの`v0.1.1` → `v0.1.2`
- `sha256`を新しい値に更新

### 5. 変更をコミット・プッシュ

```bash
git add Formula/shien.rb
git commit -m "Update shien to 0.1.2"
git push
```

## リリース後の確認

```bash
# Tapを更新
brew update

# アップグレード
brew upgrade shien

# バージョン確認
shienctl --version
```

## トラブルシューティング

### GitHub Actionsが失敗した場合

1. https://github.com/reouno/shien/actions でエラーログを確認
2. 必要に応じて`.github/workflows/release.yml`を修正
3. タグを削除して再作成：
   ```bash
   git tag -d v0.1.2
   git push origin :refs/tags/v0.1.2
   git tag -a v0.1.2 -m "Release version 0.1.2"
   git push origin v0.1.2
   ```

### Homebrewでインストールできない場合

1. Formulaの構文エラーをチェック：
   ```bash
   brew audit --strict reouno/shien/shien
   ```

2. 手動でインストールテスト：
   ```bash
   brew reinstall --build-from-source reouno/shien/shien
   ```