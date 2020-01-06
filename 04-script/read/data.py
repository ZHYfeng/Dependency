import os

from default import DependencyRPC_pb2 as pb, default


def hex_adddress(address):
    return hex(address + 0xffffffff00000000 - 5)


def uncovered_address_str(uncovered_address: pb.UncoveredAddress):
    res = ""
    res += "condition address : " + hex_adddress(uncovered_address.condition_address) + "\n"
    res += "uncovered address : " + hex_adddress(uncovered_address.uncovered_address) + "\n"
    for w in uncovered_address.write_address:
        res += "write address : " + hex_adddress(w) + "\n"
    res += "\n"
    return res


def not_covered_address_str(uncovered_address: pb.UncoveredAddress):
    res = hex_adddress(uncovered_address.condition_address) + "&" + hex_adddress(
        uncovered_address.uncovered_address) + "\n"
    return res


def not_covered_address_file_name(uncovered_address: pb.UncoveredAddress):
    res = hex_adddress(uncovered_address.condition_address) + ".txt"
    return res


def task_str(task: pb.Task):
    res = ""
    res += "*******************************************\n"
    res += "task_status : " + pb.taskStatus.Name(task.task_status) + "\n"
    res += "task priority : " + str(task.priority) + "\n"
    priority = 0
    for ua in task.uncovered_address:
        priority += task.uncovered_address[ua].priority
    res += "uncovered address priority : " + str(priority) + "\n"
    res += "condition program : " + str(task.index) + " : " + task.sig + "\n"
    res += str(task.program, default.encoding)
    res += "write address : " + hex_adddress(task.write_address) + "\n"
    res += "write program : " + str(task.write_index) + " : " + task.write_sig + "\n"
    res += str(task.write_program, default.encoding)
    res += "check_write_address : " + hex_adddress(task.check_write_address) + "\n"
    res += "-------------------------------------------\n"
    return res


def input_str(i: pb.Input):
    res = ""
    res += "*******************************************\n"
    res += "sig : " + i.sig + "\n"
    res += str(i.program, default.encoding)
    return res


class data:
    def __init__(self, dir_path):
        self.real_data = pb.Corpus()
        self.dir_path = dir_path
        self.uncovered_address_input = []
        self.uncovered_address_dependency = []

        self.read()

    def read(self):
        file_data = os.path.join(self.dir_path, default.name_data)
        print(file_data)
        if os.path.exists(file_data):
            f = open(file_data, "rb")
            self.real_data.ParseFromString(f.read())
            f.close()
            # print("self.real_data.ParseFromString(f.read())")
            self.deal()

    def deal(self):

        for a in self.real_data.uncovered_address:
            kind = self.real_data.uncovered_address[a].kind
            if kind == pb.InputRelated:
                self.uncovered_address_input.append(a)
            elif kind == pb.DependnecyRelated:
                self.uncovered_address_dependency.append(a)

        # file_result = os.path.join(self.dir_path, devices.name_data_result)
        # f = open(file_result, "w")
        # f.write(str(self.real_data))
        # f.close()

    def not_covered_address_tasks_str(self, not_covered_address):
        global ua
        res = ""
        not_covered = self.real_data.uncovered_address[not_covered_address]

        kind = 0
        # 0: default
        # 1: not find input
        # 2: not find write address
        # 3: not have write input
        # 4: cover the address
        # 5: unstable write address
        # 6: useless or FP
        # 7: unstable condition address
        # 8: unstable insert condition address

        res += "*******************************************\n"
        res += "number_arrive_basicblocks : " + str(not_covered.number_arrive_basicblocks) + "\n"
        res += "number_dominator_instructions(using) : " + str(not_covered.number_dominator_instructions) + "\n"
        res += "*******************************************\n"

        res += "# input : " + str(len(not_covered.input)) + "\n"
        for i in not_covered.input:
            if i in self.real_data.input:
                res += input_str(self.real_data.input[i])
            else:
                res += "*******************************************\n"
                res += "not find input : " + str(i) + "\n"
                kind = 1
            res += "index : " + bin(not_covered.input[i]) + "\n"
            res += "-------------------------------------------\n"

        res += "# write : " + str(len(not_covered.write_address)) + "\n"
        count = 0
        write_status = {}
        for w in self.real_data.write_address:
            write_status[w] = -1

        for w in not_covered.write_address:
            res += "## write address : " + hex(w + 0xffffffff00000000 - 5) + "\n"
            if w in self.real_data.write_address:
                write_address = self.real_data.write_address[w]
                if len(write_address.input) == 0:
                    res += "*******************************************\n"
                    res += "not have write input" + "\n"
                else:
                    count += len(write_address.input)
                    for i in write_address.input:
                        if i in self.real_data.input:
                            res += input_str(self.real_data.input[i])
                            res += "not find write input : " + str(i) + "\n"
                        res += "index : " + bin(write_address.input[i]) + "\n"
                        res += "-------------------------------------------\n"
            else:
                res += "*******************************************\n"
                res += "not find write address : " + str(w) + "\n"
                kind = 2
            res += "-------------------------------------------\n"
        if count == 0:
            kind = 3

        tasks = []
        for t in self.real_data.tasks.task_array:
            if not_covered_address in t.uncovered_address or not_covered_address in t.covered_address:
                tasks.append(t)

        res += "# tasks : " + str(len(tasks)) + "\n"
        untested_count = 0
        unstable_count = 0
        tested_count = 0
        for t in tasks:
            res += task_str(t)
            if t.task_status == pb.untested:
                untested_count += 1
                if write_status[t.write_address] < 0:
                    write_status[t.write_address] = 0
            else:
                if t.task_status == pb.tested:
                    tested_count += 1
                elif t.task_status == pb.unstable:
                    unstable_count += 1

                if not_covered_address in t.uncovered_address:
                    ua = t.uncovered_address[not_covered_address]
                    res += "task_status : " + pb.taskStatus.Name(ua.task_status) + "\n"
                    res += "check condition : " + str(ua.checkCondition) + "\n"
                    res += "chech address : " + str(ua.checkAddress) + "\n"
                    if ua.checkCondition:
                        res += ""
                    else:
                        if write_status[t.write_address] < 2:
                            write_status[t.write_address] = 2
                    res += "-------------------------------------------\n"

                    for r in t.task_run_time_data:
                        res += "-------------------------------------------\n"
                        res += "chech write address : " + str(r.check_write_address) + "\n"
                        if not_covered_address in r.uncovered_address:
                            uar = r.uncovered_address[not_covered_address]
                            res += "check condition : " + str(uar.checkCondition) + "\n"
                            res += "chech address : " + str(uar.checkAddress) + "\n"
                            if uar.checkCondition:
                                if uar.checkAddress:
                                    res += "error in ua.checkAddress" + "\n"
                                else:
                                    if r.check_write_address:
                                        res += "useless write address or FP" + "\n"
                                        if write_status[t.write_address] < 4:
                                            write_status[t.write_address] = 4
                                    else:
                                        res += "unstable insert write address" + "\n"
                                        if write_status[t.write_address] < 1:
                                            write_status[t.write_address] = 1
                            else:
                                if ua.checkCondition:
                                    res += "unstable insert condition address" + "\n"
                                    if write_status[t.write_address] < 3:
                                        write_status[t.write_address] = 3
                        res += "-------------------------------------------\n"

                elif not_covered_address in t.covered_address:
                    ua = t.covered_address[not_covered_address]
                    kind = 4

            res += "*******************************************\n"

        write_untested_count = 0
        write_unstable_count = 0
        write_useless_fp_count = 0
        condition_unstable_count = 0
        insert_condition_unstable_count = 0
        for w in write_status:
            if write_status[w] == 0:
                write_untested_count += 1
            elif write_status[w] == 1:
                write_unstable_count += 1
            elif write_status[w] == 2:
                condition_unstable_count += 1
            elif write_status[w] == 3:
                insert_condition_unstable_count += 1
            elif write_status[w] == 4:
                write_useless_fp_count += 1

        res += "untested : " + str(untested_count) + "\n"
        res += "unstable : " + str(unstable_count) + "\n"
        res += "tested : " + str(tested_count) + "\n"

        res += "write_untested_count : " + str(write_untested_count) + "\n"
        res += "write_unstable_count : " + str(write_unstable_count) + "\n"
        res += "condition_unstable_count : " + str(condition_unstable_count) + "\n"
        res += "insert condition_unstable_count : " + str(insert_condition_unstable_count) + "\n"
        res += "write_useless_fp_count : " + str(write_useless_fp_count) + "\n"
        if write_unstable_count > 0:
            kind = 5
        if condition_unstable_count > 0:
            kind = 7
        if insert_condition_unstable_count > 0:
            kind = 8
        if write_useless_fp_count > 0:
            kind = 6

        res += "kind : " + str(kind) + "\n"
        return res, kind, not_covered.number_dominator_instructions
