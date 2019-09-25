#! /usr/bin/python3
import os

import scipy.stats
import subprocess
from config import default, DependencyRPC_pb2 as pb, stats, data, axis
from config.data import uncovered_address_str, not_covered_address_str, not_covered_address_file_name


class Device:
    def __init__(self, dir_path, dir_name):
        self.dir_path = dir_path
        self.path_dev = os.path.join(self.dir_path, dir_name)

        self.file_result = os.path.join(self.path_dev, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)

        self.path_with_dra = os.path.join(self.path_dev, default.name_with_dra)
        self.results_with_dra = results(self.path_with_dra, 'C0')
        self.path_without_dra = os.path.join(self.path_dev, default.name_without_dra)
        self.results_without_dra = results(self.path_without_dra, 'C1')

        self.axises = axis.axises(self.path_dev)
        self.axises.x_axis = self.results_with_dra.axises.x_axis
        self.axises.y_axises = self.results_with_dra.axises.y_axises_statistics \
                               + self.results_without_dra.axises.y_axises_statistics
        self.axises.labels = self.results_with_dra.axises.labels_statistics \
                             + self.results_without_dra.axises.labels_statistics
        self.axises.line_styles = self.results_with_dra.axises.line_styles_statistics \
                                  + self.results_without_dra.axises.line_styles_statistics
        self.axises.colors = self.results_with_dra.axises.colors_statistics \
                             + self.results_without_dra.axises.colors_statistics

        max_coverage_with_dra = []
        for a in self.results_with_dra.axises.axises:
            max_coverage_with_dra.append(max(a.y_axis))
        max_coverage_without_dra = []
        for a in self.results_without_dra.axises.axises:
            max_coverage_without_dra.append(max(a.y_axis))
        self.statistic, self.p_value = scipy.stats.mannwhitneyu(max_coverage_with_dra, max_coverage_without_dra)

        file_figure_all = os.path.join(dir_path, dir_name, dir_name + ".pdf")
        title = " pvalue = " + str(self.p_value)
        self.axises.plot(name=file_figure_all, title=title)

        self.get_coverage()

        self.path_base = os.path.join(self.path_dev, default.name_base)
        if os.path.exists(self.path_base):
            self.basic = result(self.path_base)

            self.ca_uca_dep_with_dra = {}
            self.ca_uca_dep_without_dra = {}
            for a in self.basic.data.uncovered_address_dependency:
                if a in self.results_with_dra.max_coverage:
                    self.ca_uca_dep_with_dra[a] = self.results_with_dra.max_coverage[a]
                if a in self.results_without_dra.max_coverage:
                    self.ca_uca_dep_without_dra[a] = self.results_without_dra.max_coverage[a]

            self.ca_uca_input_with_dra = {}
            self.ca_uca_input_without_dra = {}
            for a in self.basic.data.uncovered_address_input:
                if a in self.results_with_dra.max_coverage:
                    self.ca_uca_input_with_dra[a] = self.results_with_dra.max_coverage[a]
                if a in self.results_without_dra.max_coverage:
                    self.ca_uca_input_without_dra[a] = self.results_without_dra.max_coverage[a]

            f = open(self.file_result, "a")
            f.write("=====================================================\n")
            f.write("basic:\n")
            f.write("number of uncovered address : " + str(len(self.basic.data.real_data.uncovered_address)) + "\n")
            f.write("number of uncovered address by dependency : "
                    + str(len(self.basic.data.uncovered_address_dependency)) + "\n")
            f.write("number of uncovered address by dependency covered by syzkaller with dra: "
                    + str(len(self.ca_uca_dep_with_dra)) + "\n")
            f.write("number of uncovered address by dependency covered by syzkaller without dra: "
                    + str(len(self.ca_uca_dep_without_dra)) + "\n")
            f.write(
                "number of uncovered address by input : " + str(len(self.basic.data.uncovered_address_input)) + "\n")
            f.write("number of uncovered address by input covered by syzkaller with dra: "
                    + str(len(self.ca_uca_input_with_dra)) + "\n")
            f.write("number of uncovered address by input covered by syzkaller without dra: "
                    + str(len(self.ca_uca_input_without_dra)) + "\n")

            not_covered_address_file = os.path.join(self.path_dev, "not_covered.txt")
            ff = open(not_covered_address_file, "w")
            for a in self.basic.data.uncovered_address_dependency:
                f.write(uncovered_address_str(self.basic.data.real_data.uncovered_address[a]))
                if a not in self.ca_uca_dep_with_dra and a not in self.ca_uca_dep_without_dra:
                    ff.write(not_covered_address_str(self.basic.data.real_data.uncovered_address[a]))

            ff.close()
            f.close()

            os.chdir(self.path_dev)
            cmd_a2i = default.path_a2i + " -asm=" + default.file_asm + " -objdump=" + default.file_vmlinux_objdump \
                    + " -staticRes=" + default.file_taint + " " + default.file_bc
            print(self.path_dev)
            print(cmd_a2i)
            p_a2i_img = subprocess.Popen(cmd_a2i, shell=True, preexec_fn=os.setsid)
            p_a2i_img.wait()

            for a in self.basic.data.uncovered_address_dependency:
                if a not in self.ca_uca_dep_with_dra and a not in self.ca_uca_dep_without_dra:
                    name = not_covered_address_file_name(self.basic.data.real_data.uncovered_address[a])
                    not_covered_address_file = os.path.join(self.path_dev, name)
                    ff = open(not_covered_address_file, "a")
                    ff.write(self.basic.data.not_covered_address_tasks_str(a))
                    ff.close()


        else:
            print("base not exist: " + self.path_base + "\n")

    def get_coverage(self):

        unique_coverage_with_dra = {}
        for a in self.results_with_dra.max_coverage:
            if a not in self.results_without_dra.max_coverage:
                unique_coverage_with_dra[a] = self.results_with_dra.max_coverage[a]

        unique_coverage_without_dra = {}
        for a in self.results_without_dra.max_coverage:
            if a not in self.results_with_dra.max_coverage:
                unique_coverage_without_dra[a] = self.results_without_dra.max_coverage[a]

        max_coverage = {}
        for a in self.results_without_dra.max_coverage:
            max_coverage[a] = self.results_without_dra.max_coverage[a]
        for a in self.results_with_dra.max_coverage:
            if a not in max_coverage:
                max_coverage[a] = 0
            else:
                max_coverage[a] = max_coverage[a] + self.results_with_dra.max_coverage[a]

        f = open(self.file_result, "a")
        f.write("=====================================================\n")
        f.write("coverage:\n")
        f.write("unique_coverage_with_dra : " + str(len(unique_coverage_with_dra)) + "\n")
        f.write(str(unique_coverage_with_dra) + "\n")
        f.write("unique_coverage_without_dra : " + str(len(unique_coverage_without_dra)) + "\n")
        f.write(str(unique_coverage_without_dra) + "\n")
        f.write("max_coverage : " + str(len(max_coverage)) + "\n")
        f.close()


class results:
    def __init__(self, dir_path, color=''):
        self.dir_path = dir_path
        self.color = color

        self.file_result = os.path.join(self.dir_path, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)

        self.results = []
        self.statistics = stats.stats(self.dir_path)
        self.axises = axis.axises(self.dir_path, self.color)

        self.max_coverage = {}
        self.uncovered_address_input = []
        self.uncovered_address_dependency = []
        self.max_uncoverage = {}

        self.deal_results()
        self.get_uncovered_address()
        self.get_max_coverage()

    def deal_results(self):
        if os.path.exists(self.dir_path):
            for (dir_path, dir_names, file_names) in os.walk(self.dir_path):
                for file_name in file_names:
                    if file_name.startswith(default.name_stat):
                        r = result(dir_path)
                        self.results.append(r)
                        self.statistics.statistics.append(r.stat)
                        self.axises.axises.append(r.axis)

        self.axises.deal()

        f = open(self.file_result, "a")
        f.write("=====================================================\n")
        f.write("stat:\n")
        self.statistics.get_average()
        f.write(str(self.statistics.processed_stat.real_stat))
        self.statistics.processed_stat.deal_stat()
        f.write(str(self.statistics.processed_stat.real_stat))
        f.close()

    def get_uncovered_address(self):
        for r in self.results:
            for a in r.data.real_data.uncovered_address:
                self.max_uncoverage[a] = r.data.real_data.uncovered_address[a]
                kind = r.data.real_data.uncovered_address[a].kind
                if kind == pb.InputRelated:
                    self.uncovered_address_input.append(a)
                elif kind == pb.DependnecyRelated:
                    self.uncovered_address_dependency.append(a)

        remove_address = []
        for a in self.max_uncoverage:
            for r in self.results:
                if a not in r.data.real_data.uncovered_address:
                    remove_address.append(a)
                    break

        for a in remove_address:
            self.max_uncoverage.pop(a)

        f = open(self.file_result, "a")
        f.write("=====================================================\n")
        f.write("uncovered address: " + str(len(self.max_uncoverage)) + "\n")
        for a in self.max_uncoverage:
            f.write(uncovered_address_str(self.max_uncoverage[a]))
        f.close()

    def get_max_coverage(self):
        for s in self.statistics.statistics:
            for a in s.real_stat.coverage.coverage:
                if a not in self.max_coverage:
                    self.max_coverage[a] = 0
                else:
                    self.max_coverage[a] = self.max_coverage[a] + 1


class result:
    def __init__(self, dir_path):
        self.dir_path = dir_path
        self.file_result = os.path.join(self.dir_path, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)
        self.data = data.data(self.dir_path)
        self.stat = stats.stat(self.dir_path)
        self.stat.get_time_coverage()
        self.axis = axis.axis(self.dir_path, self.stat.x_axis, self.stat.y_axis, '-')


def get_stat_file(path):
    is_dev = False
    is_results = False
    is_result = False

    dir_name = os.path.basename(path)
    dir_path = os.path.dirname(path)
    if dir_name.startswith(default.name_dev):
        Device(dir_path, dir_name)
    elif dir_name.startswith(default.name_with_dra) or dir_name.startswith(default.name_without_dra):
        path_results = os.path.join(dir_path, dir_name)
        results(path_results)
    elif dir_name.startswith(default.name_stat):
        result(dir_path)
    else:
        for (dir_path, dir_names, file_names) in os.walk(path):
            for dir_name in dir_names:
                if dir_name.startswith(default.name_dev):
                    is_dev = True
                    Device(dir_path, dir_name)
            if is_dev:
                break

            for dir_name in dir_names:
                if dir_name.startswith(default.name_with_dra) or dir_name.startswith(default.name_without_dra):
                    is_dev = True
                    path_results = os.path.join(dir_path, dir_name)
                    results(path_results)
            if is_results:
                break

            for file_name in file_names:
                if file_name.startswith(default.name_stat):
                    is_result = True
                    result(dir_path)
            if is_result:
                break
