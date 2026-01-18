# Release Guide

For 0.X versions (pre-1.0).

## Quick Release

```bash
git status
make ci-check

git tag -a v0.1.0 -m "Release v0.1.0

Features:
- One-liner install  
- GitHub auth
- SSH key management"

git push origin v0.1.0
```

CI/CD automatically builds and uploads binaries.

## Manual Build

```bash
make ci-build
ls -lh dist/
cd dist && sha256sum * > CHECKSUMS.txt
```

## Hotfix

```bash
git checkout -b hotfix/description
# ... fix ...
make test
git push origin hotfix/description
gh pr create --title "HOTFIX: description"

# After merge:
git tag -a v0.1.1 -m "Hotfix: description"
git push origin v0.1.1
```

---

**Repository**: https://github.com/IniZio/nexus
