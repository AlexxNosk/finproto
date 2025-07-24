#!/bin/bash
set -e

TMP_DIR="/tmp/finam-trade-api"
PROTO_SRC_DIR="$TMP_DIR/proto/grpc/tradeapi/v1"
PROTO_DST_DIR="./trade_api/grpc/tradeapi/v1"
GOOGLEAPIS_DIR="./trade_api/googleapis"
OUT_DIR="."

echo "📥 Cloning finam-trade-api..."
rm -rf "$TMP_DIR"
git clone --depth 1 git@github.com:FinamWeb/finam-trade-api.git "$TMP_DIR"

echo "📁 Copying .proto files to $PROTO_DST_DIR"
mkdir -p "$PROTO_DST_DIR"
cp -r "$PROTO_SRC_DIR"/* "$PROTO_DST_DIR/"

echo "🧹 Cleaning up temp"
rm -rf "$TMP_DIR"

# ✅ Check if googleapis is installed
if [ ! -d "$GOOGLEAPIS_DIR" ]; then
    echo "📦 Cloning googleapis into $GOOGLEAPIS_DIR"
    mkdir -p "$(dirname "$GOOGLEAPIS_DIR")"
    git clone --depth 1 https://github.com/googleapis/googleapis.git "$GOOGLEAPIS_DIR"
else
    echo "✅ googleapis already exists at $GOOGLEAPIS_DIR"
fi

# Find all .proto files and compile them
echo "🔍 Finding .proto files..."
find "$PROTO_DST_DIR" -name "*.proto" | while read -r proto_file; do
    echo "📦 Compiling $proto_file..."
    protoc \
        --proto_path=. \
        --proto_path=trade_api \
        --proto_path="$GOOGLEAPIS_DIR" \
        --go_out="$OUT_DIR" \
        --go-grpc_out="$OUT_DIR" \
        "$proto_file"
done

echo "✅ Done."
