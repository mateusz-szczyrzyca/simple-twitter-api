#!/usr/bin/env python3

import subprocess

output = ""
command = "doctl compute droplet list --no-header --format Name,PublicIPv4,PrivateIPv4,Tags | grep -v -E 'S4|abstynenci.pl'"
cmd = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
for line in cmd.stdout.readlines():
    output += line.decode("utf-8")

num = 0
tags = {}
format = ""
for line in output.strip().split("\n"):
    line = " ".join(line.split())
    l = 0
    public_host = ""
    private_host = ""
    for hostdata in line.split(" "):
        if l == 0:
            format += "["+hostdata+"]\n"
        if l == 1:
            format += hostdata+" "
            public_host = hostdata
        if l == 2:
            format += "private_ip="+hostdata+" "
            private_host = hostdata
        if l == 3 and hostdata != "":
            format += "tag="+hostdata+" "
            pubkey = hostdata + "_public"
            prvkey = hostdata + "_private"
            if pubkey not in tags:
                tags[pubkey] = []
            if prvkey not in tags:
                tags[prvkey] = []
            tags[pubkey].append(public_host)
            tags[prvkey].append(private_host)

        l = l+1
    format += "\n\n"

for k in tags.keys():
    format += "["+k+"]\n"
    format += "\n".join(tags[k])
    format += "\n\n"

print(format)

