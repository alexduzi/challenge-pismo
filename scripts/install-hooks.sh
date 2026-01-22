#!/bin/bash

# Install git hooks for the project
# Run this script after cloning the repository: ./scripts/install-hooks.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "🔧 Installing git hooks..."
echo ""

# Check if we're in a git repository
if [ ! -d "$PROJECT_ROOT/.git" ]; then
    echo "❌ Error: Not a git repository. Please run this from the project root."
    exit 1
fi

# Create pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash

# Git pre-commit hook
# Runs unit tests before allowing commit

echo "🔍 Running pre-commit checks..."
echo ""

# Run unit tests
echo "🧪 Running unit tests..."
if ! make test-unit; then
    echo ""
    echo "❌ Unit tests failed! Please fix the failing tests before committing."
    exit 1
fi

echo "✅ Unit tests passed!"
echo ""
echo "✨ All pre-commit checks passed! Proceeding with commit..."
echo ""

exit 0
EOF

# Make the hook executable
chmod +x "$HOOKS_DIR/pre-commit"

echo "✅ Git hooks installed successfully!"
echo ""
echo "Installed hooks:"
echo "  - pre-commit: Runs unit tests before commit"
echo ""
echo "To skip hooks temporarily, use: git commit --no-verify"