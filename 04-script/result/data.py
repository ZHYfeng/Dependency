import os
from result import DependencyRPC_pb2 as pb
from result import default


class data:
    def __init__(self, dir_path):
        self.data = pb.Corpus()
        self.dir_path = dir_path

    def data_read(self):
        file_data = os.path.join(self.dir_path, default.name_data)
        if os.path.exists(file_data):
            f = open(file_data, "rb")
            self.data.ParseFromString(f.read())
            f.close()
            # data_deal(dir_path, data)

    def data_deal(self):
        file_result = os.path.join(self.dir_path, default.name_data_result)
        f = open(file_result, "w")
        f.write(str(self.data))
        f.close()
