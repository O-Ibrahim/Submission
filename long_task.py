import time

i = 0
while True:
    i += 1
    time.sleep(0.01)
    print(i % 999999999999999999999, flush=True)
