import os

import Statistic_pb2 as pb, default


class stats:
    def __init__(self, dir_path):
        self.dir_path = dir_path
        self.statistics = []
        self.processed_stat = stat()

    def read(self):
        if os.path.exists(self.dir_path):
            for (dir_path, dir_names, file_names) in os.walk(self.dir_path):
                for file_name in file_names:
                    if file_name.startswith(default.name_stat):
                        s = stat(dir_path)
                        self.statistics.append(s)

class stat:
    def __init__(self, dir_path=""):
        self.dir_path = dir_path
        self.file_result = os.path.join(self.dir_path, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)
        self.real_stat = pb.Statistics()
        self.processed_stat = pb.Statistics()
        self.x_axis = []
        self.y_axis = []
        self.read()

    def read(self):
        file_stat = os.path.join(self.dir_path, default.name_stat)
        if os.path.exists(file_stat):
            f = open(file_stat, "rb")
            self.real_stat.ParseFromString(f.read())
            f.close()

    def get_time_coverage(self):
        t0 = 0
        num = 0
        for i in self.real_stat.coverage.time:
            while i.time > t0:
                t0 = t0 + 60
                self.x_axis.append(t0)
                if i.time > t0:
                    self.y_axis.append(num)
                else:
                    num = i.num
                    self.y_axis.append(num)