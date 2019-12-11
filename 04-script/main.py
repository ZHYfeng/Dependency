import sys

from config import config
from default import default
from read import result

if __name__ == "__main__":
    if len(sys.argv) > 1:
        if sys.argv[1] == "generate":
            if len(sys.argv) > 2:
                default.default.path_result = sys.argv[2]
            config.generate_dev_dir()
        if sys.argv[1] == "read" and len(sys.argv) > 2:
            if len(sys.argv) > 3:
                default.do_figure = False
            result.read_results(sys.argv[2])
