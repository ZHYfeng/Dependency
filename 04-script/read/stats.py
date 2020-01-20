import os

from default import Statistic_pb2 as pb, Base_pb2, default


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

    def get_average(self):
        for ss in self.statistics:
            for s in ss.real_stat.stat:
                self.processed_stat.real_stat.stat[s].name = ss.real_stat.stat[s].name
                s_add(self.processed_stat.real_stat.stat[s], ss.real_stat.stat[s])
        s_length = len(self.statistics)
        for s in self.processed_stat.real_stat.stat:
            self.processed_stat.real_stat.stat[s].executeNum = int(
                self.processed_stat.real_stat.stat[s].executeNum / s_length)
            self.processed_stat.real_stat.stat[s].time = int(self.processed_stat.real_stat.stat[s].time / s_length)
            self.processed_stat.real_stat.stat[s].newTestCaseNum = int(
                self.processed_stat.real_stat.stat[s].newTestCaseNum / s_length)
            self.processed_stat.real_stat.stat[s].newAddressNum = int(
                self.processed_stat.real_stat.stat[s].newAddressNum / s_length)


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
            self.deal_stat()

            # file_result = os.path.join(self.dir_path, devices.name_stat_result)
            # f = open(file_result, "w")
            # f.write(str(self.real_stat))
            # f.write(str(self.deal_stat))
            # f.close()

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

    def deal_stat(self):
        s_copy(self.processed_stat.stat[pb.StatGenerate], self.real_stat.stat[pb.StatGenerate])
        s_copy(self.processed_stat.stat[pb.StatFuzz], self.real_stat.stat[pb.StatFuzz])
        s_copy(self.processed_stat.stat[pb.StatCandidate], self.real_stat.stat[pb.StatCandidate])
        s_copy(self.processed_stat.stat[pb.StatTriage], self.real_stat.stat[pb.StatTriage])
        s_add(self.processed_stat.stat[pb.StatTriage], self.real_stat.stat[pb.StatMinimize])
        s_copy(self.processed_stat.stat[pb.StatSmash], self.real_stat.stat[pb.StatSmash])
        s_add(self.processed_stat.stat[pb.StatSmash], self.real_stat.stat[pb.StatHint])
        s_add(self.processed_stat.stat[pb.StatSmash], self.real_stat.stat[pb.StatSeed])
        s_copy(self.processed_stat.stat[pb.StatDependency], self.real_stat.stat[pb.StatDependency])
        s_copy(self.processed_stat.stat[pb.StatDependency], self.real_stat.stat[pb.StatDependencyBoot])

        if self.dir_path != "":
            f = open(self.file_result, "a")
            f.write("=====================================================\n")
            f.write("useful input 1`43v b:" + "\n")
            for ui in self.real_stat.useful_input:
                if ui.num > 100:
                    f.write("program : " + str(ui.input.program) + "\n")
                    f.write("number : " + str(ui.num) + "\n")
                    f.write("time : " + str(ui.time) + "\n")
                    f.write("address : " + str(ui.new_address) + "\n")
            f.close()


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
