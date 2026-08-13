#!/usr/bin/env bash
# check-plugin-isolation.sh
# Đảm bảo nguyên tắc plugin-like: pkg/budget và pkg/rules KHÔNG bị import vào file lõi bị cấm.
# Chạy trước khi commit: bash scripts/check-plugin-isolation.sh
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

# Các file lõi bị cấm import feature mới
FORBIDDEN_FILES=(
  "pkg/models/transaction.go"
  "pkg/models/transaction_category.go"
  "pkg/services/transactions.go"
  "pkg/services/transaction_categories.go"
)

# Các import prefix bị cấm trong file lõi
FORBIDDEN_IMPORTS=(
  "github.com/mayswind/ezbookkeeping/pkg/budget"
  "github.com/mayswind/ezbookkeeping/pkg/rules"
)

EXIT_CODE=0

for file in "${FORBIDDEN_FILES[@]}"; do
  if [[ ! -f "$file" ]]; then
    continue
  fi
  for imp in "${FORBIDDEN_IMPORTS[@]}"; do
    # Tìm dòng import chứa forbidden prefix
    if grep -nE "^[[:space:]]*\"${imp}" "$file" >/dev/null 2>&1; then
      echo "❌ VIOLATION: '$file' imports '$imp'"
      grep -nE "^[[:space:]]*\"${imp}" "$file"
      EXIT_CODE=1
    fi
  done
done

if [[ $EXIT_CODE -eq 0 ]]; then
  echo "✅ Plugin isolation OK: no forbidden imports in core files."
fi

exit $EXIT_CODE
