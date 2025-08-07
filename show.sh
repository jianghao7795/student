git log --all --pretty=format: --name-only | sort | uniq | while read f; do 
  if file -b --mime-type "$f" 2>/dev/null | grep -q '^application/octet-stream$'; then 
    echo "Binary file found: $f"; 
  fi; 
done
