import os
from result import DependencyRPC_pb2 as pb
from result import default


class statistics:
    def __init__(self):
        self.statistics = []

    def read(self, path):
        if os.path.exists(path):
            for (dir_path, dir_names, file_names) in os.walk(path):
                for file_name in file_names:
                    if file_name.startswith(default.name_stat):
                        s = statistic(dir_path)
                        s.read()
                        self.statistics.append(s)

    def get_average(self):
        r = pb.Statistics()
        for stat in self.statistics:
            for s in stat.stat:
                r.stat[s].name = stat.stat[s].name
                s_add(r.stat[s], stat.stat[s])
        s_length = len(self.statistics)
        for s in r.stat:
            r.stat[s].executeNum = int(r.stat[s].executeNum / s_length)
            r.stat[s].time = int(r.stat[s].time / s_length)
            r.stat[s].newTestCaseNum = int(r.stat[s].newTestCaseNum / s_length)
            r.stat[s].newAddressNum = int(r.stat[s].newAddressNum / s_length)
        return r


class statistic:
    def __init__(self, dir_path):
        self.dir_path = dir_path
        self.stat = pb.Statistics()
        self.deal = pb.Statistics()
        self.read()

    def read(self):
        file_stat = os.path.join(self.dir_path, default.name_stat)
        f = open(file_stat, "rb")
        self.stat.ParseFromString(f.read())
        f.close()
        file_result = os.path.join(self.dir_path, default.name_stat_result)
        self.f = open(file_result, "w")
        self.out_put(str(self.stat))
        self.out_put(self.deal(self.stat))
        f.close()

    def out_put(self, context):
        self.f.write(context)

    def get_time_coverage(self):
        x_axis = []
        y_axis = []
        t0 = 0
        num = 0
        for i in self.stat.coverage.time:
            while i.time > t0:
                t0 = t0 + 60
                x_axis.append(t0)
                if i.time > t0:
                    y_axis.append(num)
                else:
                    num = i.num
                    y_axis.append(num)

        return x_axis, y_axis

    def stat_deal(self):
        self.deal = pb.Statistics()
        s_copy(self.deal.stat[pb.StatGenerate], self.stat.stat[pb.StatGenerate])
        s_copy(self.deal.stat[pb.StatFuzz], self.stat.stat[pb.StatFuzz])
        s_copy(self.deal.stat[pb.StatCandidate], self.stat.stat[pb.StatCandidate])
        s_copy(self.deal.stat[pb.StatTriage], self.stat.stat[pb.StatTriage])
        s_add(self.deal.stat[pb.StatTriage], self.stat.stat[pb.StatMinimize])
        s_copy(self.deal.stat[pb.StatSmash], self.stat.stat[pb.StatSmash])
        s_add(self.deal.stat[pb.StatSmash], self.stat.stat[pb.StatHint])
        s_add(self.deal.stat[pb.StatSmash], self.stat.stat[pb.StatSeed])
        s_copy(self.deal.stat[pb.StatDependency], self.stat.stat[pb.StatDependency])


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
