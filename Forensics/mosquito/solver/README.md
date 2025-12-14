dump elf:
`tshark -r mosquito.pcap -Y "(tcp) && (frame.len == 584)" -T fields -e data -e tcp.segment_data | xxd -r -p > tes.elf`

ct filter:
`(_ws.col.protocol == "MQTT") && (mqtt.msgtype == 3)`

dec:
```
#!/usr/bin/env python3
import hashlib, binascii
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

TS = 1761146087
USERNAME  = "0367d08072f248a0474784700be2b224084581469eb00ef5e827e7d1782e34ff"
PASSWORD  = "26852622c1f9c54ea24190e5bb33790c7b5442c2c79c4509309f1c2d468fe384"
CT_HEX    = "eedb69d9d5bda1593f9ef278b5cefa0f05334f80b7d8d90b75c15d6cf8090770eb766465177a56ac5557aab8af520bc4"

def s256_bytes(b: bytes) -> bytes: return hashlib.sha256(b).digest()

K_text = s256_bytes(USERNAME.encode())[:16]
IV_text = s256_bytes(f"{PASSWORD}:{TS}".encode())[:16]

K_hex = s256_bytes(binascii.unhexlify(USERNAME))[:16]
IV_hex = s256_bytes(binascii.unhexlify(PASSWORD) + b":" + str(TS).encode())[:16]

print("Go-matching  KEY:", K_text.hex())
print("Go-matching  IV :", IV_text.hex())

ct = bytes.fromhex(CT_HEX)
for label, K, IV in [
    ("Go-matching TEXT inputs", K_text, IV_text),
    ("Hex-decoded inputs (alt)", K_hex, IV_hex),
]:
    try:
        pt = unpad(AES.new(K, AES.MODE_CBC, IV).decrypt(ct), 16)
        print(f"[OK] {label}: {pt!r}")
    except Exception as e:
        print(f"[!!] {label}: padding/decrypt failed ({e})")
```