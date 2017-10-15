import json
import sys

filename = sys.argv[1]

with open(filename) as f:
    records = json.load(f)

for r in records:
    s = ""
    s += r['QB']['name'] + ","
    s += r['RB1']['name'] + ","
    s += r['RB2']['name'] + ","
    s += r['WR1']['name'] + ","
    s += r['WR2']['name'] + ","
    s += r['WR3']['name'] + ","
    s += r['TE']['name'] + ","
    s += r['FLEX']['name'] + ","
    s += r['DST']['name'] + ","
    s += str(r['Spent'])
    print s
