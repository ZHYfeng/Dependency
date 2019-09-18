import os

from result import DependencyRPC_pb2 as pb
from result import default


def uncovered_address_str(uncovered_address: pb.UncoveredAddress):
    res = ""
    res += "condition address : " + hex(uncovered_address.condition_address + 0xffffffff00000000 - 5) + "\n"
    res += "uncovered address : " + hex(uncovered_address.uncovered_address + 0xffffffff00000000 - 5) + "\n"
    for w in uncovered_address.write_address:
        res += "write address : " + hex(w + 0xffffffff00000000 - 5) + "\n"
    res += "\n"
    return res


class data:
    def __init__(self, dir_path):
        self.data = pb.Corpus()
        self.dir_path = dir_path

    def read(self):
        file_data = os.path.join(self.dir_path, default.name_data)
        if os.path.exists(file_data):
            f = open(file_data, "rb")
            self.data.ParseFromString(f.read())
            f.close()
            # data_deal(dir_path, data)

    def deal(self):
        file_result = os.path.join(self.dir_path, default.name_data_result)
        f = open(file_result, "w")
        f.write(str(self.data))
        f.close()
