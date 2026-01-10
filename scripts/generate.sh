#!/bin/bash
set -e

echo "🧹 Cleaning previous generated files..."

# Remove previous generated proto files
rm -rf server/proto-generated/*
rm -rf client/src/proto-generated/*

echo "🔨 Generating code from proto files..."

# Generate Go and TypeScript code
buf generate

echo "🔨 Generating Ent code..."

# Generate Ent code from schema
cd server && GOWORK=off go generate ./ent && cd ..

echo "✅ Code generation completed!"
