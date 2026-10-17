#!/usr/bin/env bash
# Regenerates the self-hosted Kairos brand faces under web-common/static/fonts/kairos.
#
# We serve these ourselves rather than linking the Google Fonts CDN: the CDN sees
# every visitor's IP, and it costs two third-party origins in the CSP. Run this
# to pick up upstream font revisions or to add a subset.
set -euo pipefail

# Inter Tight and Source Sans 3 are the brand; the rest belong to optional
# Styles (Geist: nova, glass; Source Serif 4: editorial; Atkinson: lumen;
# Spectral, Karla, Poppins: apetic). A face downloads only when a page uses it.
FAMILIES="family=Inter+Tight:wght@500..700&family=Source+Sans+3:wght@400..600&family=Geist:wght@400..600&family=Geist+Mono:wght@400..600&family=Source+Serif+4:opsz,wght@8..60,500..600&family=Atkinson+Hyperlegible+Next:wght@400..700&family=Spectral:wght@600;700&family=Karla:wght@400..700&family=Poppins:wght@500;600"
# Spanish needs only latin; latin-ext covers the rest of western Europe. Adding
# more subsets here is the supported way to widen coverage.
SUBSETS="latin latin-ext"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${ROOT}/web-common/static/fonts/kairos"
OUT_CSS="${ROOT}/web-common/static/fonts/kairos.css"

# Google serves woff2 only to user agents known to support it.
UA="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

mkdir -p "${OUT_DIR}"
source_css="$(mktemp)"
trap 'rm -f "${source_css}"' EXIT

curl -sS -A "${UA}" "https://fonts.googleapis.com/css2?${FAMILIES}&display=swap" -o "${source_css}"

python3 - "${source_css}" "${OUT_DIR}" "${OUT_CSS}" "${SUBSETS}" <<'PY'
import pathlib
import re
import subprocess
import sys

source_css, out_dir, out_css, subsets = sys.argv[1:5]
wanted = set(subsets.split())
blocks = re.findall(
    r"/\* ([a-z-]+) \*/\s*(@font-face \{.*?\})",
    pathlib.Path(source_css).read_text(),
    re.S,
)

faces = []
for subset, block in blocks:
    if subset not in wanted:
        continue
    family = re.search(r"font-family:\s*'([^']+)'", block).group(1)
    url = re.search(r"url\((https://[^)]+)\)", block).group(1)
    # A static family (no weight range) ships one file per weight.
    weight = re.search(r"font-weight:\s*([^;]+);", block).group(1).strip()
    slug = family.lower().replace(" ", "-")
    name = f"{slug}-{subset}.woff2" if " " in weight else f"{slug}-{weight}-{subset}.woff2"
    subprocess.run(["curl", "-sS", "-o", f"{out_dir}/{name}", url], check=True)
    faces.append(block.replace(url, f"/fonts/kairos/{name}").replace("'", '"'))

header = """/* Kairos brand faces, served from our own origin.
 *
 * Self-hosted rather than pulled from the Google Fonts CDN: the CDN sees every
 * visitor's IP, which sits badly with the data residency this product sells, and
 * it needs two third-party origins in the CSP that ADR-0006 set out to remove.
 * Same-origin also means html-to-image can inline these faces into a PDF export
 * under connect-src 'self' (see features/exports/pdf/capture.ts).
 *
 * Generated file, do not edit. Regenerate with scripts/fetch-brand-fonts.sh
 */

"""
pathlib.Path(out_css).write_text(header + "\n\n".join(faces) + "\n")
print(f"wrote {len(faces)} faces to {out_css}")
PY
