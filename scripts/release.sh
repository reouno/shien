#!/bin/bash
set -e

# 色付き出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 使用方法
if [ $# -ne 1 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 0.1.2"
    exit 1
fi

VERSION=$1
TAG="v${VERSION}"

echo -e "${GREEN}Starting release process for version ${VERSION}...${NC}"

# 1. 現在のブランチを確認
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}Error: You must be on the main branch to release${NC}"
    exit 1
fi

# 2. 変更がコミットされているか確認
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes${NC}"
    exit 1
fi

# 3. タグを作成してプッシュ
echo -e "${YELLOW}Creating and pushing tag ${TAG}...${NC}"
git tag -a "$TAG" -m "Release version ${VERSION}"
git push origin "$TAG"

# 4. GitHub Actionsの完了を待つ
echo -e "${YELLOW}Waiting for GitHub Actions to complete the release...${NC}"
echo "Please check: https://github.com/reouno/shien/actions"
echo "Press Enter when the release is complete..."
read -r

# 5. リリースファイルのSHA256を取得
echo -e "${YELLOW}Downloading release and calculating SHA256...${NC}"
RELEASE_URL="https://github.com/reouno/shien/releases/download/${TAG}/shien-darwin-arm64.tar.gz"
curl -sL "$RELEASE_URL" -o /tmp/shien-release.tar.gz
SHA256=$(shasum -a 256 /tmp/shien-release.tar.gz | awk '{print $1}')
echo "SHA256: $SHA256"

# 6. homebrew-shienリポジトリの場所を確認
HOMEBREW_REPO_PATH=""
if [ -d "$HOME/homebrew-shien" ]; then
    HOMEBREW_REPO_PATH="$HOME/homebrew-shien"
elif [ -d "../homebrew-shien" ]; then
    HOMEBREW_REPO_PATH="../homebrew-shien"
else
    echo -e "${YELLOW}Enter the path to homebrew-shien repository:${NC}"
    read -r HOMEBREW_REPO_PATH
fi

# 7. Formula を更新
echo -e "${YELLOW}Updating Homebrew formula...${NC}"
FORMULA_PATH="$HOMEBREW_REPO_PATH/Formula/shien.rb"
sed -i.bak "s/version \".*\"/version \"${VERSION}\"/" "$FORMULA_PATH"
sed -i.bak "s|download/v[0-9.]\+/|download/${TAG}/|" "$FORMULA_PATH"
sed -i.bak "s/sha256 \".*\"/sha256 \"${SHA256}\"/" "$FORMULA_PATH"
rm "$FORMULA_PATH.bak"

# 8. Formulaの変更をコミット・プッシュ
cd "$HOMEBREW_REPO_PATH"
git add Formula/shien.rb
git commit -m "Update shien to ${VERSION}"
git push

echo -e "${GREEN}Release completed successfully!${NC}"
echo -e "${GREEN}Users can now update with: brew upgrade shien${NC}"

# クリーンアップ
rm -f /tmp/shien-release.tar.gz