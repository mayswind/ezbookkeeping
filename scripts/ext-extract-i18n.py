#!/usr/bin/env python3
"""Adds every label used in src/ext to src/ext/locales/en.json.

Like the app's own en.json, the English sentence is the key. Run it after adding a label:
    python3 scripts/ext-extract-i18n.py
Existing messages (and the nested "ext" messages with placeholders) are never overwritten.
Other languages: add src/ext/locales/<language code>.json with the same keys; missing keys fall back to English.
"""
import json
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
SRC = ROOT / "src" / "ext"
EN = SRC / "locales" / "en.json"
LITERAL = re.compile(r"""\btt\(\s*(['"])((?:\\.|(?!\1).)*)\1""")


def labels():
    found = {}
    for path in sorted(list(SRC.rglob("*.vue")) + list(SRC.rglob("*.ts"))):
        if "__tests__" in path.parts:
            continue
        for match in LITERAL.finditer(path.read_text()):
            text = match.group(2).replace("\\'", "'").replace('\\"', '"')
            if not text.startswith("ext."):  # keys with a dot prefix are nested messages, kept by hand
                found.setdefault(text, str(path.relative_to(ROOT)))
    return found


def main():
    current = json.loads(EN.read_text()) if EN.exists() else {}
    added = [text for text in labels() if text not in current]
    for text in added:
        current[text] = text
    ordered = {k: current[k] for k in sorted((k for k in current if k != "ext"), key=str.lower)}
    if "ext" in current:
        ordered["ext"] = current["ext"]
    EN.write_text(json.dumps(ordered, indent=4, ensure_ascii=False) + "\n")
    print(f"{len(added)} label(s) added, {len(ordered)} entries in {EN.relative_to(ROOT)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
