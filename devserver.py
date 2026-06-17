#!/usr/bin/env python3
"""
Dev server for LukasDerBaum website.
Watches source files, rebuilds on change, and live-reloads the browser.

Usage:
    python devserver.py [--port 8080] [--dir .] [--build-cmd "go run ."]
"""

import argparse
import functools
import hashlib
import http.server
import os
import queue
import subprocess
import sys
import threading
import time
from pathlib import Path

# ── ANSI colours ─────────────────────────────────────────────────────────────
GREEN  = "\033[92m"
YELLOW = "\033[93m"
RED    = "\033[91m"
CYAN   = "\033[96m"
RESET  = "\033[0m"

def log(colour, tag, msg):
    ts = time.strftime("%H:%M:%S")
    print(f"{colour}[{ts}] [{tag}]{RESET} {msg}")

# ── File-hash watcher ─────────────────────────────────────────────────────────
WATCH_EXTENSIONS = {".go", ".templ", ".md", ".css", ".html", ".toml", ".yaml",".ts"}
IGNORE_DIRS      = {"build", ".git", "vendor", "__pycache__", "_templ.go"}

def file_hash(path: Path) -> str:
    try:
        return hashlib.md5(path.read_bytes()).hexdigest()
    except OSError:
        return ""

def scan(root: Path) -> dict[str, str]:
    hashes = {}
    for dirpath, dirs, files in os.walk(root):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
        for fname in files:
            p = Path(dirpath) / fname
            if p.suffix in WATCH_EXTENSIONS:
                if str(p).endswith("_templ.go"):
                    continue
                hashes[str(p.relative_to(root))] = file_hash(p)
    return hashes

# ── SSE broadcast hub ─────────────────────────────────────────────────────────
class Hub:
    def __init__(self):
        self._clients: list[queue.Queue] = []
        self._lock = threading.Lock()

    def add(self, q: queue.Queue):
        with self._lock:
            self._clients.append(q)

    def remove(self, q: queue.Queue):
        with self._lock:
            try:
                self._clients.remove(q)
            except ValueError:
                pass

    def reload(self):
        with self._lock:
            dead = []
            for q in self._clients:
                try:
                    q.put_nowait("reload")
                except Exception:
                    dead.append(q)
            for q in dead:
                self._clients.remove(q)
        log(CYAN, "SSE", f"Sent reload to {len(self._clients)} client(s)")

hub = Hub()

# ── Live-reload script injected into every HTML response ──────────────────────
LIVE_RELOAD_SCRIPT = b"""
<script>
(function(){
  const es = new EventSource('/__livereload');
  es.onmessage = () => { console.log('[devserver] reloading'); location.reload(); };
  es.onerror   = () => console.warn('[devserver] SSE connection lost, will retry');
})();
</script>
"""

# ── HTTP handler ──────────────────────────────────────────────────────────────
class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, build_dir: str = "build", **kwargs):
        self.build_dir = build_dir
        super().__init__(*args, directory=build_dir, **kwargs)

    def do_GET(self):
        if self.path == "/__livereload":
            self._handle_sse()
            return

        # Resolve path to a file on disk
        fs_path = self.translate_path(self.path)
        if os.path.isdir(fs_path):
            fs_path = os.path.join(fs_path, "index.html")

        if fs_path.endswith(".html") and os.path.isfile(fs_path):
            with open(fs_path, "rb") as f:
                raw = f.read()

            # Inject reload script before </body>, or append if missing
            if b"</body>" in raw:
                patched = raw.replace(b"</body>", LIVE_RELOAD_SCRIPT + b"</body>", 1)
            else:
                patched = raw + LIVE_RELOAD_SCRIPT

            self.send_response(200)
            self.send_header("Content-Type", "text/html; charset=utf-8")
            self.send_header("Content-Length", str(len(patched)))
            self.send_header("Cache-Control", "no-cache")
            self.end_headers()
            self.wfile.write(patched)
        else:
            super().do_GET()

    def _handle_sse(self):
        self.send_response(200)
        self.send_header("Content-Type", "text/event-stream")
        self.send_header("Cache-Control", "no-cache")
        self.send_header("Connection", "keep-alive")
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()

        q: queue.Queue = queue.Queue()
        hub.add(q)
        try:
            while True:
                try:
                    msg = q.get(timeout=25)
                    self.wfile.write(f"data: {msg}\n\n".encode())
                    self.wfile.flush()
                except queue.Empty:
                    # heartbeat to keep connection alive
                    self.wfile.write(b": ping\n\n")
                    self.wfile.flush()
        except (BrokenPipeError, ConnectionResetError):
            pass
        finally:
            hub.remove(q)

    def log_message(self, fmt, *args):
        pass  # suppress access log noise

# ── Build runner ──────────────────────────────────────────────────────────────
def run_build(cmd: str, cwd: Path) -> bool:
    log(YELLOW, "BUILD", f"Running: {cmd}")
    t0 = time.monotonic()
    result = subprocess.run(cmd, shell=True, cwd=cwd, capture_output=True, text=True)
    elapsed = time.monotonic() - t0
    if result.returncode == 0:
        log(GREEN, "BUILD", f"OK ({elapsed:.2f}s)")
        if result.stdout.strip():
            print(result.stdout.strip())
        return True
    else:
        log(RED, "BUILD", f"FAILED ({elapsed:.2f}s)")
        print(result.stderr.strip() or result.stdout.strip())
        return False

# ── Watcher loop ──────────────────────────────────────────────────────────────
def watcher(root: Path, build_cmd: str, interval: float = 0.5):
    log(CYAN, "WATCH", f"Watching {root}  (extensions: {', '.join(sorted(WATCH_EXTENSIONS))})")
    prev = scan(root)

    while True:
        time.sleep(interval)
        curr = scan(root)

        added   = set(curr) - set(prev)
        removed = set(prev) - set(curr)
        changed = {k for k in curr if k in prev and curr[k] != prev[k]}

        if added | removed | changed:
            for f in sorted(added):   log(GREEN,  "NEW",     f)
            for f in sorted(removed): log(RED,    "DELETED", f)
            for f in sorted(changed): log(YELLOW, "CHANGED", f)

            if run_build(build_cmd, root):
                hub.reload()

        prev = curr

# ── Entry point ───────────────────────────────────────────────────────────────
def main():
    parser = argparse.ArgumentParser(description="Hot-reload dev server")
    parser.add_argument("--port",      type=int,   default=8080,            help="HTTP port (default 8080)")
    parser.add_argument("--dir",       type=str,   default=".",             help="Project root (default: cwd)")
    parser.add_argument("--build-cmd", type=str,   default="make build_dev",help='Build command (default: "make build_dev")')
    parser.add_argument("--build-dir", type=str,   default="build",         help='Dir to serve (default: "build")')
    parser.add_argument("--interval",  type=float, default=0.5,             help="Poll interval in seconds (default 0.5)")
    args = parser.parse_args()

    root      = Path(args.dir).resolve()
    build_dir = root / args.build_dir

    log(CYAN, "START", "Initial build…")
    run_build(args.build_cmd, root)

    if not build_dir.exists():
        log(RED, "ERROR", f"Build dir '{build_dir}' does not exist after build. Check your build command.")
        sys.exit(1)

    t = threading.Thread(target=watcher, args=(root, args.build_cmd, args.interval), daemon=True)
    t.start()

    handler = functools.partial(Handler, build_dir=str(build_dir))
    server  = http.server.ThreadingHTTPServer(("", args.port), handler)

    log(GREEN, "SERVER", f"Serving http://localhost:{args.port}  (build dir: {build_dir})")
    log(CYAN,  "INFO",   "Press Ctrl+C to stop")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        log(YELLOW, "STOP", "Shutting down")

if __name__ == "__main__":
    main()