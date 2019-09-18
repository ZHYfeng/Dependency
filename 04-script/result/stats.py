import os
from result import DependencyRPC_pb2 as pb
from result import default


class stats:
    def __init__(self, dir_path):
        self.dir_path = dir_path
        self.statistics = []
        self.stat = stat()

    def read(self):
        if os.path.exists(self.dir_path):
            for (dir_path, dir_names, file_names) in os.walk(self.dir_path):
                for file_name in file_names:
                    if file_name.startswith(default.name_stat):
                        s = stat(dir_path)
                        s.read()
                        self.statistics.append(s)

    def get_average(self):
        for ss in self.statistics:
            for s in ss.stat:
                self.stat.stat.stat[s].name = ss.stat[s].name
                s_add(self.stat.stat[s], ss.stat[s])
        s_length = len(self.statistics)
        for s in self.stat.stat:
            self.stat.stat[s].executeNum = int(self.stat.stat[s].executeNum / s_length)
            self.stat.stat[s].time = int(self.stat.stat[s].time / s_length)
            self.stat.stat[s].newTestCaseNum = int(self.stat.stat[s].newTestCaseNum / s_length)
            self.stat.stat[s].newAddressNum = int(self.stat.stat[s].newAddressNum / s_length)


class stat:
    def __init__(self, dir_path=""):
        self.dir_path = dir_path
        self.stat = pb.Statistics()
        self.stat_deal = pb.Statistics()
        self.x_axis = []
        self.y_axis = []

    def read(self):
        file_stat = os.path.join(self.dir_path, default.name_stat)
        if os.path.exists(file_stat):
            f = open(file_stat, "rb")
            self.stat.ParseFromString(f.read())
            f.close()

            # file_result = os.path.join(self.dir_path, default.name_stat_result)
            # f = open(file_result, "w")
            # f.write(str(self.stat))
            # self.deal()
            # f.write(str(self.stat_deal))
            # f.close()

    def get_time_coverage(self):
        t0 = 0
        num = 0
        for i in self.stat.coverage.time:
            while i.time > t0:
                t0 = t0 + 60
                self.x_axis.append(t0)
                if i.time > t0:
                    self.y_axis.append(num)
                else:
                    num = i.num
                    self.y_axis.append(num)

    def deal(self):
        s_copy(self.stat_deal.stat[pb.StatGenerate], self.stat.stat[pb.StatGenerate])
        s_copy(self.stat_deal.stat[pb.StatFuzz], self.stat.stat[pb.StatFuzz])
        s_copy(self.stat_deal.stat[pb.StatCandidate], self.stat.stat[pb.StatCandidate])
        s_copy(self.stat_deal.stat[pb.StatTriage], self.stat.stat[pb.StatTriage])
        s_add(self.stat_deal.stat[pb.StatTriage], self.stat.stat[pb.StatMinimize])
        s_copy(self.stat_deal.stat[pb.StatSmash], self.stat.stat[pb.StatSmash])
        s_add(self.stat_deal.stat[pb.StatSmash], self.stat.stat[pb.StatHint])
        s_add(self.stat_deal.stat[pb.StatSmash], self.stat.stat[pb.StatSeed])
        s_copy(self.stat_deal.stat[pb.StatDependency], self.stat.stat[pb.StatDependency])


def s_add(stat1: pb.Statistic, stat2: pb.Statistic) -> pb.Statistic:
    stat1.executeNum += stat2.executeNum
    stat1.time += stat2.time
    stat1.newTestCaseNum += stat2.newTestCaseNum
    stat1.newAddressNum += stat2.newAddressNum
    return stat1


def s_copy(stat1: pb.Statistic, stat2: pb.Statistic) -> pb.Statistic:
    stat1.name = stat2.name
    stat1.executeNum = stat2.executeNum
    stat1.time = stat2.time
    stat1.newTestCaseNum = stat2.newTestCaseNum
    stat1.newAddressNum = stat2.newAddressNum
    return stat1
