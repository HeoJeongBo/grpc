#!/bin/bash
set -e

echo "🔨 Generating code from proto files..."

# Generate Go and TypeScript code
buf generate

echo "✅ Code generation completed!"
echo ""
echo "Generated files:"
echo "  - Go: server/proto-generated/item/v1/*.go"
echo "  - TypeScript: client/src/proto-generated/item/v1/*.ts"
