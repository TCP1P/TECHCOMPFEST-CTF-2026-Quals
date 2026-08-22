import randcrack
import random
import os
from libnum import s2n
import signal

used = set()
seed = os.urandom(12)
wack = randcrack.RandCrack()
random.seed(seed)

print("Teio step!")
for _ in range(624):
    signal.alarm(2) # Dont waste time
    num = int(input(">> "))
    signal.alarm(0)

    if num in used:
        print("No duplicates!")
        exit(1)
    used.add(num)

    wack.submit(num)
    print(random.getrandbits(32))

sta = 0
for _ in range(624):
    sta ^= random.getrandbits(32) ^ wack.predict_getrandbits(32)

if abs(sta) < s2n(b'bakushinbakushin') % 175433:
    print("Winning live complete!")
    gift = open('flag.txt').read()
    print(gift)
else:
    print("!")
    exit(1)