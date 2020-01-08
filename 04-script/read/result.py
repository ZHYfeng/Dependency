#! /usr/bin/python3
import os
import subprocess

import scipy.stats

from default import DependencyRPC_pb2 as pb, default
from read import data, stats, axis
from read.data import uncovered_address_str, not_covered_address_str, not_covered_address_file_name, hex_adddress


class Device:
    def __init__(self, dir_path, dir_name):
        self.unique_coverage_without_dra = {}
        self.unique_coverage_with_dra = {}
        self.dir_path = dir_path
        self.dir_name = dir_name
        self.path_dev = os.path.join(self.dir_path, dir_name)

        self.file_result = os.path.join(self.path_dev, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)

        self.path_with_dra = os.path.join(self.path_dev, default.name_with_dra)
        self.results_with_dra = Results(self.path_with_dra, 'C0')
        self.path_without_dra = os.path.join(self.path_dev, default.name_without_dra)
        self.results_without_dra = Results(self.path_without_dra, 'C1')

        print(self.path_dev)
        self.axises = axis.axises(self.path_dev)
        self.set_axises()

        self.statistic, self.p_value = 0, 0
        self.get_mann_withney_utest()

        self.get_coverage()
        self.get_base()

    def set_axises(self):
        self.axises.x_axis = self.results_with_dra.axises.x_axis
        self.axises.y_axises = self.results_with_dra.axises.y_axises_statistics + self.results_without_dra.axises.y_axises_statistics
        self.axises.labels = self.results_with_dra.axises.labels_statistics + self.results_without_dra.axises.labels_statistics
        self.axises.line_styles = self.results_with_dra.axises.line_styles_statistics + self.results_without_dra.axises.line_styles_statistics
        self.axises.colors = self.results_with_dra.axises.colors_statistics + self.results_without_dra.axises.colors_statistics

    def get_mann_withney_utest(self):
        max_coverage_with_dra = []
        for a in self.results_with_dra.axises.axises:
            max_coverage_with_dra.append(max(a.y_axis))
        max_coverage_without_dra = []
        for a in self.results_without_dra.axises.axises:
            max_coverage_without_dra.append(max(a.y_axis))
        self.statistic, self.p_value = scipy.stats.mannwhitneyu(
            max_coverage_with_dra, max_coverage_without_dra)

        file_figure_all = os.path.join(self.dir_path, self.dir_name, self.dir_name + ".pdf")
        title = "pvalue = " + str(self.p_value)
        self.axises.plot(name=file_figure_all, title=title)

    def get_coverage(self):
        for a in self.results_with_dra.max_coverage:
            if a not in self.results_without_dra.max_coverage:
                self.unique_coverage_with_dra[a] = self.results_with_dra.max_coverage[a]

        for a in self.results_without_dra.max_coverage:
            if a not in self.results_with_dra.max_coverage:
                self.unique_coverage_without_dra[a] = self.results_without_dra.max_coverage[a]

        union_coverage = {}
        for a in self.results_without_dra.max_coverage:
            union_coverage[a] = self.results_without_dra.max_coverage[a]
        for a in self.results_with_dra.max_coverage:
            if a not in union_coverage:
                union_coverage[a] = self.results_with_dra.max_coverage[a]
            else:
                union_coverage[a] = union_coverage[a] + self.results_with_dra.max_coverage[a]

        intersection_coverage = {}
        for a in union_coverage:
            if a in self.results_with_dra.max_coverage and a in self.results_without_dra.max_coverage:
                intersection_coverage[a] = self.results_without_dra.max_coverage[a] + \
                                           self.results_with_dra.max_coverage[a]

        f = open(self.file_result, "a")
        f.write("=====================================================\n")
        f.write("coverage:\n")
        f.write("unique_coverage_with_dra : " +
                str(len(self.unique_coverage_with_dra)) + "\n")
        f.write(str(self.unique_coverage_with_dra) + "\n")

        f.write("unique_coverage_without_dra : " +
                str(len(self.unique_coverage_without_dra)) + "\n")
        f.write(str(self.unique_coverage_without_dra) + "\n")
        f.write("union_coverage : " + str(len(union_coverage)) + "\n")
        f.write("intersection_coverage : " + str(len(intersection_coverage)) + "\n")

        solved_condition = {}
        for r in self.results_with_dra.results:
            for t in r.data.real_data.tasks.task_array:
                for ca in t.covered_address:
                    solved_condition[ca] = t.covered_address[ca]

        stable_solved_conditon = {}
        unstable_solved_conditon = {}
        for a in solved_condition:
            if a in self.results_with_dra.max_coverage:
                stable_solved_conditon[a] = solved_condition[a]
            else:
                unstable_solved_conditon[a] = solved_condition[a]

        f.write("solved_condition : " + str(len(solved_condition)) + "\n")
        f.write("stable_solved_conditon : " + str(len(stable_solved_conditon)) + "\n")
        f.write("unstable_solved_conditon : " + str(len(unstable_solved_conditon)) + "\n")

        f.close()

    def get_base(self):
        path_base = os.path.join(self.path_dev, default.name_base)
        if os.path.exists(path_base):
            basic = result(path_base)

            f = open(self.file_result, "a")
            f.write("=====================================================\n")
            f.write("basic:\n")

            ca_uca_dep_with_dra = {}
            ca_uca_dep_without_dra = {}
            for a in basic.data.uncovered_address_dependency:
                if a in self.results_with_dra.max_coverage:
                    ca_uca_dep_with_dra[a] = self.results_with_dra.max_coverage[a]
                if a in self.results_without_dra.max_coverage:
                    ca_uca_dep_without_dra[a] = self.results_without_dra.max_coverage[a]
            f.write("number of uncovered address : " + str(len(basic.data.real_data.uncovered_address)) + "\n")
            f.write("related to dependency: " + str(len(basic.data.uncovered_address_dependency)) + "\n")
            count = 0
            for a in ca_uca_dep_with_dra:
                count += ca_uca_dep_with_dra[a]
            f.write("covered by syzkaller with dra: "
                    + str(len(ca_uca_dep_with_dra)) + " count : " + str(count) + "\n")

            addresses_dependency_covered_with_dra_dependency = {}
            for a in ca_uca_dep_with_dra:
                for s in self.results_with_dra.statistics.statistics:
                    if a in s.real_stat.coverage.coverage:
                        if s.real_stat.coverage.coverage[a] == 1:
                            addresses_dependency_covered_with_dra_dependency[a] = 1
                            break
            f.write(
                "covered by dependency mutate: " + str(len(addresses_dependency_covered_with_dra_dependency)) + "\n")

            count = 0
            for a in addresses_dependency_covered_with_dra_dependency:
                if a in self.unique_coverage_with_dra:
                    count = count + 1
            f.write("unique one of them: " + str(count) + "\n")

            count = 0
            for a in ca_uca_dep_without_dra:
                count += ca_uca_dep_without_dra[a]
            f.write("covered by syzkaller without dra: " + str(len(ca_uca_dep_without_dra)) + " count : " + str(
                count) + "\n")

            # ca_uca_input_with_dra = {}
            # ca_uca_input_without_dra = {}
            # for a in basic.data.uncovered_address_input:
            #     if a in self.results_with_dra.max_coverage:
            #         ca_uca_input_with_dra[a] = self.results_with_dra.max_coverage[a]
            #     if a in self.results_without_dra.max_coverage:
            #         ca_uca_input_without_dra[a] = self.results_without_dra.max_coverage[a]
            # f.write("number of uncovered address by input : " + str(len(basic.data.uncovered_address_input)) + "\n")
            # f.write("number of uncovered address by input covered by syzkaller with dra: "
            #         + str(len(ca_uca_input_with_dra)) + "\n")
            # f.write("number of uncovered address by input covered by syzkaller without dra: "
            #         + str(len(ca_uca_input_without_dra)) + "\n")

            not_covered_address_file = os.path.join(self.path_dev, "not_covered.txt")
            ff = open(not_covered_address_file, "w")
            for a in basic.data.uncovered_address_dependency:
                f.write(uncovered_address_str(basic.data.real_data.uncovered_address[a]))
                if a not in self.results_with_dra.max_coverage and a not in self.results_without_dra.max_coverage:
                    ff.write(not_covered_address_str(basic.data.real_data.uncovered_address[a]))

            ff.close()

            os.chdir(self.path_dev)

            cmd_rm_0x = "rm -rf 0x*"
            p_rm_0x = subprocess.Popen(cmd_rm_0x, shell=True, preexec_fn=os.setsid)
            p_rm_0x.wait()
            # cmd_a2i = default.path_a2i + " -asm=" + default.file_asm + " -objdump=" + default.file_vmlinux_objdump \
            #           + " -staticRes=./" + default.file_taint + " -function=./" + \
            #           default.file_function + " " + default.file_bc
            # print(cmd_a2i)
            # p_a2i_img = subprocess.Popen(cmd_a2i, shell=True, preexec_fn=os.setsid)
            # p_a2i_img.wait()

            ua_status = {}
            ua_insts = {}
            ua_input = {}
            ua_write = {}
            ua_count = {}
            ua_tasks = {}
            ua_tested_tasks = {}
            for a in basic.data.uncovered_address_dependency:
                if a in self.results_with_dra.uncovered_address_dependency and \
                        a in self.results_without_dra.uncovered_address_dependency:

                    # for a in self.results_with_dra.uncovered_address_dependency:
                    name = not_covered_address_file_name(basic.data.real_data.uncovered_address[a])
                    not_covered_address_file = os.path.join(self.path_dev, name)
                    print(not_covered_address_file)
                    ff = open(not_covered_address_file, "a")
                    for r in self.results_with_dra.results:
                        if a in r.data.real_data.uncovered_address:
                            res, kind, inst, input_count, write_count, count, tasks, tested_tasks = r.data.not_covered_address_tasks_str(
                                a)
                            ff.write(res)
                            ua_status[a] = kind
                            ua_insts[a] = inst
                            ua_input[a] = input_count
                            ua_write[a] = write_count
                            ua_count[a] = count
                            ua_tasks[a] = tasks
                            ua_tested_tasks[a] = tested_tasks
                            break
                    ff.close()
            ua_status_count = {x: 0 for x in range(9)}
            for ua in ua_status:
                ua_status_count[ua_status[ua]] += 1

            res = ""
            res += "untested : " + str(ua_status_count[0]) + "\n"
            res += "not find input : " + str(ua_status_count[1]) + "\n"
            res += "not find write address : " + str(ua_status_count[2]) + "\n"
            res += "not have write input : " + str(ua_status_count[3]) + "\n"
            res += "cover : " + str(ua_status_count[4]) + "\n"
            res += "unstable write : " + str(ua_status_count[5]) + "\n"
            res += "unstable condition : " + str(ua_status_count[7]) + "\n"
            res += "unstable insert condition : " + str(ua_status_count[8]) + "\n"
            res += "useless or FP : " + str(ua_status_count[6]) + "\n"
            f.write(res)

            res = ""
            sort_ua_insts = sorted(ua_insts.items(), key=lambda kv: kv[1])
            for ua in sort_ua_insts:
                res += "uncovered address : " + hex_adddress(ua[0]) + " inst : " + str(ua[1]).zfill(4) + " kind : " \
                       + str(ua_status[ua[0]]) + " input : " + str(ua_input[ua[0]]).zfill(3) + " write : " + str(
                    ua_write[ua[0]]).zfill(3) + " count : " + str(ua_count[ua[0]]).zfill(7) + " tasks : " + str(
                    ua_tasks[ua[0]]).zfill(5) + " tested tasks : " + str(ua_tested_tasks[ua[0]]).zfill(4) + "\n"

            res += "\n"

            cover_addres = {}
            task_priority = {}
            for r in self.results_with_dra.results:
                for t in r.data.real_data.tasks.task_map:
                    tt = r.data.real_data.tasks.task_map[t]
                    priority = 0
                    for ua in tt.uncovered_address:
                        priority += tt.uncovered_address[ua].priority
                    priority = priority * pow(2, tt.priority)
                    task_priority[t] = priority
                    for ca in tt.covered_address:
                        cover_addres[ca] = t.covered_address[ca]

            res += "cover_addres in task : " + str(len(cover_addres)) + "\n"

            sort_task_priority = sorted(task_priority.items(), key=lambda kv: kv[1])
            for r in self.results_with_dra.results:
                for t in sort_task_priority:
                    tt = r.data.real_data.tasks.task_map[t[0]]
                    res += "task : " + " final priority : " + str(t[1]).zfill(16) + " priority : " + str(tt.priority).zfill(4) + " task status : " + str(pb.taskStatus.Name(
                        tt.task_status)).zfill(8) + " uncovered address : " + str(
                        len(tt.uncovered_address)).zfill(3) + " execute count : " + str(tt.count).zfill(3) + "\n"

            f.write(res)
            f.close()

        else:
            print("base not exist: " + path_base + "\n")


class Results:
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
                    if file_name == default.name_stat:
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
                    self.max_coverage[a] = 1
                else:
                    self.max_coverage[a] = self.max_coverage[a] + 1


class result:
    def __init__(self, dir_path):
        self.dir_path = dir_path
        self.file_result = os.path.join(self.dir_path, default.name_data_result)
        if os.path.exists(self.file_result):
            os.remove(self.file_result)
        # print("self.data = data.data(self.dir_path)")
        self.data = data.data(self.dir_path)
        # print("self.stat = stats.stat(self.dir_path)")
        self.stat = stats.stat(self.dir_path)
        # print("self.stat.get_time_coverage()")
        self.stat.get_time_coverage()
        # print("self.axis = axis.axis(self.dir_path, self.stat.x_axis, self.stat.y_axis, '-')")
        self.axis = axis.axis(self.dir_path, self.stat.x_axis, self.stat.y_axis, '-')


def read_results(path):
    is_dev = False
    is_results = False
    is_result = False

    dir_name = os.path.basename(path)
    dir_path = os.path.dirname(path)
    if dir_name.startswith(default.name_dev):
        Device(dir_path, dir_name)
    elif dir_name.startswith(default.name_with_dra) or dir_name.startswith(default.name_without_dra):
        path_results = os.path.join(dir_path, dir_name)
        Results(path_results)
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
                    Results(path_results)
            if is_results:
                break

            for file_name in file_names:
                if file_name.startswith(default.name_stat):
                    is_result = True
                    result(dir_path)
            if is_result:
                break
