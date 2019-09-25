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
    res += "priority : " + task.priority
    res += "condition program : " + task.index + " : " + task.sig + "\n"
    res += task.program
    res += "write address : " + task.write_address + "\n"
    res += "write program : " + task.write_index + " : " + task.write_sig + "\n"
    res += task.write_program
    res += "check_write_address : " + task.check_write_address + "\n"
    res += "check_write_address_final : " + task.check_write_address_final + "\n"
    res += "check_write_address_remove : " + task.check_write_address_remove + "\n"
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

        tasks = []
        for t in self.real_data.tasks:
            if not_covered_address in t.uncovered_address:
                tasks.append(t)

        res = "# tasks : " + str(len(tasks)) + "\n"
        for t in tasks:
            res += task_str(t)
        return res
