#!/bin/bash
set -e

echo "🔨 Generating code from proto files..."

# Generate Go and TypeScript code
buf generate

echo "🔨 Generating Ent code..."

# Generate Ent code from schema
cd server && GOWORK=off go generate ./ent && cd ..

echo "✅ Code generation completed!"
echo ""
echo "Generated files:"
echo "  - Go: server/proto-generated/item/v1/*.go"
echo "  - TypeScript: client/src/proto-generated/item/v1/*.ts"
echo "  - Ent: server/ent/*.go"
