import os

from config import default, DependencyRPC_pb2 as pb


def uncovered_address_str(uncovered_address: pb.UncoveredAddress):
    res = ""
    res += "condition address : " + hex(uncovered_address.condition_address + 0xffffffff00000000 - 5) + "\n"
    res += "uncovered address : " + hex(uncovered_address.uncovered_address + 0xffffffff00000000 - 5) + "\n"
    for w in uncovered_address.write_address:
        res += "write address : " + hex(w + 0xffffffff00000000 - 5) + "\n"
    res += "\n"
    return res


def not_covered_address_str(uncovered_address: pb.UncoveredAddress):
    res = hex(uncovered_address.condition_address + 0xffffffff00000000 - 5) + "&" + hex(
        uncovered_address.uncovered_address + 0xffffffff00000000 - 5) + "\n"
    return res


def not_covered_address_file_name(uncovered_address: pb.UncoveredAddress):
    res = hex(uncovered_address.condition_address + 0xffffffff00000000 - 5) + ".txt"
    return res


def task_str(task: pb.Task):
    res = ""
    res += "-------------------------------------------\n"
    res += "task_status : " + str(task.task_status) + "\n"
    res += "priority : " + str(task.priority) + "\n"
    res += "condition program : " + str(task.index) + " : " + task.sig + "\n"
    res += str(task.program, default.encoding)
    res += "write address : " + str(task.write_address) + "\n"
    res += "write program : " + str(task.write_index) + " : " + task.write_sig + "\n"
    res += str(task.write_program, default.encoding)
    res += "check_write_address : " + str(task.check_write_address) + "\n"
    res += "check_write_address_final : " + str(task.check_write_address_final) + "\n"
    res += "check_write_address_remove : " + str(task.check_write_address_remove) + "\n"
    res += "-------------------------------------------\n"
    return res


def input_str(i : pb.Input):
    res = ""
    res += "-------------------------------------------\n"
    res += str(i.program, default.encoding)
    res += "-------------------------------------------\n"
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
        if os.path.exists(file_data):
            f = open(file_data, "rb")
            self.real_data.ParseFromString(f.read())
            f.close()

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
        res = ""
        not_covered = self.real_data.uncovered_address[not_covered_address]

        res += "# input : " + str(len(not_covered.input)) + "\n"
        for i in not_covered.input:
            if i in self.real_data.input:
                res += input_str(self.real_data.input[i])
            else:
                res += "-------------------------------------------\n"
                res += "not find input : " + str(i) + "\n"
            res += "index : " + bin(not_covered.input[i]) + "\n"
            res += "-------------------------------------------\n"

        res += "# write : " + str(len(not_covered.write_address)) + "\n"
        for w in not_covered.write_address:
            res += "## write address : " + hex(w + 0xffffffff00000000 - 5) + "\n"
            if w in self.real_data.write_address:
                write_address = self.real_data.write_address[w]
                for i in write_address.input:
                    if i in self.real_data.input:
                        res += input_str(self.real_data.input[i])
                    else:
                        res += "-------------------------------------------\n"
                        res += "not find write input : " + str(i) + "\n"
                    res += "index : " + bin(write_address.input[i]) + "\n"
                    res += "-------------------------------------------\n"
            else:
                res += "-------------------------------------------\n"
                res += "not find write address : " + str(w) + "\n"
            res += "-------------------------------------------\n"

        tasks = []
        for t in self.real_data.tasks.task:
            if not_covered_address in t.uncovered_address:
                tasks.append(t)

        res += "# tasks : " + str(len(tasks)) + "\n"
        for t in tasks:
            res += task_str(t)
        return res
