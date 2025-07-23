#!/bin/bash
# 识别并删除二进制文件
git ls-tree -r --name-only HEAD | while read filename; do
  if file -b --mime-type "$filename" | grep -q '^application/octet-stream$'; then
    git filter-branch --force --index-filter \
    "git rm --cached --ignore-unmatch '$filename'" \
    --prune-empty --tag-name-filter cat -- --all
  fi
done
