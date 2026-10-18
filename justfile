# Tasks for the website; run them in `nix develop` (see flake.nix).
#
# The pages are built from a checkout of arkadia-lang/arkadia: its guide,
# its changelog and, through its ark command, the documentation of its
# standard library. Point `arkadia` at it (`just arkadia=path site`).

arkadia := "../arkadia"

# List the recipes.
default:
    @just --list

# Build the website into _site/.
site:
    #!/usr/bin/env bash
    set -euo pipefail
    arkadia=$(cd {{arkadia}} && pwd)
    (cd "$arkadia" && just build)
    rm -rf _site
    mkdir -p _site/assets
    cp style.css _site/
    cp "$arkadia/docs/assets/logo.png" _site/assets/
    page() { pandoc --template template.html -t html5 -s "$@"; }
    page -f markdown+east_asian_line_breaks -V nav-home=true -o _site/index.html index.md
    page -f gfm+east_asian_line_breaks --lua-filter guide.lua --toc --toc-depth=2 \
        -M pagetitle=言語ガイド -V nav-guide=true -o _site/guide.html "$arkadia/docs/guide.md"
    page -f gfm+east_asian_line_breaks -M pagetitle=変更履歴 -V nav-changelog=true \
        -o _site/changelog.html "$arkadia/CHANGELOG.md"
    ARKADIA_ROOT=$arkadia go run ./stdref -ark "$arkadia/bin/ark" -lib "$arkadia/lib" -out _site/std
    # the commit of Arkadia the site was built from (see .github/workflows/pages.yml)
    git -C "$arkadia" rev-parse HEAD > _site/arkadia-commit.txt
    echo "built _site/"

# Format the Go sources.
fmt:
    gofmt -w stdref
