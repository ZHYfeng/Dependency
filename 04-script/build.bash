#!/bin/bash
echo "[*] Trying to build python"
sudo apt install -y python3-pip python3-tk
python3 -m pip install --user matplotlib
python3 -m pip install --user python3-protobuf
python3 -m pip install --user protobuf
python3 -m pip install --user scipy