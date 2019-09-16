import sys
import argparse

from config import config
from result import read_result

if __name__ == "__main__":
    if len(sys.argv) > 1:
        if sys.argv[1] == "generate":
            if len(sys.argv) > 2:
                config.devices.path_result = sys.argv[2]
            config.generate_dev_dir()
        if sys.argv[1] == "read" and len(sys.argv) > 2:
            if len(sys.argv) > 3:
                read_result.do_figure = False
            read_result.get_stat_file(sys.argv[2])

