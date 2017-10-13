import csv

def normalize(name):
    content = []
    with open("./csv/" + name, "r") as f:
        content = f.readlines()
        f.close()

    output = ""
    prevLine = ""
    for i, line in enumerate(content):
        if i % 2 == 0:
            prevLine = line.strip()
        else:
            prevLine += line
            output += prevLine

    with open("./csv/" + name, "w") as f:
        f.write(output)
        f.close()
