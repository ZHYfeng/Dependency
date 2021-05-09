import sys

import default, config

if __name__ == "__main__":
    if len(sys.argv) > 1:
        if sys.argv[1] == "generate":
            if len(sys.argv) > 2:
                default.path_result = sys.argv[2]
            config.generate_dev_dir()
