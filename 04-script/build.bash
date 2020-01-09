#!/bin/bash
echo "[*] Trying to build python"

python3 -m pip install --user matplotlib
python3 -m pip install --user python3-protobuf
python3 -m pip install --user protobuf
python3 -m pip install --user scipy