from ecdsa import ellipticcurve
from libnum import n2s, s2n
from random import randint
from ecdsa.numbertheory import square_root_mod_prime

def lift_x(curve: ellipticcurve.CurveFp, x: int):
    a, b, p = curve.a(), curve.b(), curve.p()
    rhs = (pow(x, 3, p) + a * x + b) % p
    
    # 2. Try to find the square root
    try:
        y = square_root_mod_prime(rhs, p)
        return ellipticcurve.Point(curve, x, y)
    except:
        return None

def get_random_point(curve: ellipticcurve.CurveFp):
    P = None
    x = randint(1, curve.p() - 1)
    while not P:
        P = lift_x(curve, x)
        x = randint(1, curve.p() - 1)
    return P

def print_point(P: ellipticcurve.Point):
    print(f"({P.x()}, {P.y()})")

p = 94124950304330993908925892037528579476182856287778334781204189865278579685821
a = 23021404719642270564990343636058205767306310438190499409503140948031340113729
b = 16276897285419258155793642760334397668199916080103049064087720015882865744187

flag = open("flag.txt", "rb").read()
flag = int(s2n(flag))

print("Starting up...")
ec = ellipticcurve.CurveFp(p, a, b)

for _ in range(100):
    print("What you want?")
    choice = int(input("> "))
    
    if choice == 1:
        print(f"Printing {flag.bit_length()} lines...")
        while flag:
            b = flag & 1
            flag >>= 1
            y = get_random_point(ec)
            x = get_random_point(ec)
            c = y * (p >> 202) * randint(202, 20202) + x * b
            print_point(c)
    if choice == 2:
        try:
            H = lift_x(ec, flag) * 6969
            print(True)
        except:
            print(randint(0, 1) == 1)
    if choice == 3:
        d = int(input("d?> "))
        flag += d