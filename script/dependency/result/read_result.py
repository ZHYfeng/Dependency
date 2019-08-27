#! /usr/bin/python3
import os
import sys

import matplotlib.pyplot as plt

sys.path.append(os.getcwd())
from script.dependency.dra import DependencyRPC_pb2 as pb

do_figure = True

name_dev = "dev_"
name_with_dra = "result-with-dra"
name_without_dra = "result-without-dra"
name_stat = "statistics.bin"
name_stat_result = "statistics.txt"
name_data = "data.bin"
name_data_result = "data.txt"
length = 1 * 24 * 60
time_run = length * 60  # second


def stat_print(f, stat):
    f.write(str(stat))


def stat_read(dir_path):
    file_stat = os.path.join(dir_path, name_stat)
    stat = pb.Statistics()
    f = open(file_stat, "rb")
    stat.ParseFromString(f.read())
    f.close()
    file_result = os.path.join(dir_path, name_stat_result)
    f = open(file_result, "w")
    stat_print(f, stat)
    stat_print(f, stat_deal(stat))
    f.close()

    return stat


def data_read(dir_path):
    file_data = os.path.join(dir_path, name_data)
    data = pb.Corpus()
    if os.path.exists(file_data):
        f = open(file_data, "rb")
        data.ParseFromString(f.read())
        f.close()

        data_deal(dir_path, data)
    return data


def data_deal(dir_path, data):
    file_result = os.path.join(dir_path, name_data_result)
    f = open(file_result, "w")
    f.write(str(data))
    f.close()


def stats_read(path):
    stats = []
    if os.path.exists(path):
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.startswith(name_stat):
                    stats.append(stat_read(dir_path))

    return stats


def axis_plot(name, x_axis, y_axis):
    if not do_figure or len(y_axis) == 0:
        return
    f = plt.figure()
    plt.plot(x_axis, y_axis)
    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    f.savefig(fname=name, bbox_inches='tight', format="pdf")
    plt.close(f)


def axises_plot(name, x_axis, y_axises, labels):
    if not do_figure or len(y_axises) == 0:
        return
    f = plt.figure()
    for i in range(len(labels)):
        plt.plot(x_axis, y_axises[i], label=labels[i])

    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    if len(labels) != 0:
        plt.legend()
    f.savefig(fname=name, bbox_inches='tight', format="pdf")
    plt.close(f)


def stat_get_time_coverage(stat):
    x_axis = []
    y_axis = []
    t0 = 0
    num = 0
    for i in stat.coverage.time:
        while i.time > t0:
            t0 = t0 + 60
            x_axis.append(t0)
            if i.time > t0:
                y_axis.append(num)
            else:
                num = i.num
                y_axis.append(num)

    return x_axis, y_axis


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


def stat_deal(stat):
    r = pb.Statistics()
    s_copy(r.stat[pb.StatGenerate], stat.stat[pb.StatGenerate])
    s_copy(r.stat[pb.StatFuzz], stat.stat[pb.StatFuzz])
    s_copy(r.stat[pb.StatCandidate], stat.stat[pb.StatCandidate])
    s_copy(r.stat[pb.StatTriage], stat.stat[pb.StatTriage])
    s_add(r.stat[pb.StatTriage], stat.stat[pb.StatMinimize])
    s_copy(r.stat[pb.StatSmash], stat.stat[pb.StatSmash])
    s_add(r.stat[pb.StatSmash], stat.stat[pb.StatHint])
    s_add(r.stat[pb.StatSmash], stat.stat[pb.StatSeed])
    s_copy(r.stat[pb.StatDependency], stat.stat[pb.StatDependency])
    return r


def stats_get_average(stats):
    r = pb.Statistics()
    for stat in stats:
        for s in stat.stat:
            r.stat[s].name = stat.stat[s].name
            s_add(r.stat[s], stat.stat[s])
    s_length = len(stats)
    for s in r.stat:
        r.stat[s].executeNum = int(r.stat[s].executeNum / s_length)
        r.stat[s].time = int(r.stat[s].time / s_length)
        r.stat[s].newTestCaseNum = int(r.stat[s].newTestCaseNum / s_length)
        r.stat[s].newAddressNum = int(r.stat[s].newAddressNum / s_length)
    return r


def result_deal(dir_path):
    stat = stat_read(dir_path)
    x_axis, y_axis = stat_get_time_coverage(stat)
    file_figure = os.path.join(dir_path, "coverage.pdf")
    axis_plot(file_figure, x_axis, y_axis)

    data = data_read(dir_path)
    return stat, x_axis, y_axis, data


def axis_expansion(length, x_axises, y_axises):
    for x in x_axises:
        if len(x) != 0:
            max_time = x[-1]
        else:
            max_time = 0
        for i in range(length - len(x)):
            max_time = max_time + 60
            x.append(max_time)

    for y in y_axises:
        if len(y) != 0:
            max_num = y[-1]
        else:
            max_num = 0
        for i in range(length - len(y)):
            y.append(max_num)
    return x_axises, y_axises


def result_get_from_sub_dir(stats, x_axises, y_axises, labels, path, name_label):
    if os.path.exists(path):
        s, x, y = results_deal(path)
        stats.append(s)
        x_axises.append(x)
        y_axises.append(y)
        labels.append(name_label)


def results_deal(path):
    stats = []
    x_axises = []
    y_axises = []
    data = []
    labels = []
    if os.path.exists(path):
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.startswith(name_stat):
                    s, x, y, d = result_deal(dir_path)
                    stats.append(s)
                    x_axises.append(x)
                    y_axises.append(y)
                    data.append(d)
                    labels.append(dir_path)

    stat = stats_get_average(stats)
    file_result = os.path.join(path, name_stat_result)
    f = open(file_result, "w")
    stat_print(f, stat)
    stat_print(f, stat_deal(stat))
    f.close()

    x_axises, y_axises = axis_expansion(length, x_axises, y_axises)
    x_axis = [sum(e) / len(e) for e in zip(*x_axises)]
    y_axis = [sum(e) / len(e) for e in zip(*y_axises)]
    if len(y_axises) == 0:
        return stats, x_axis, y_axis
    file_figure_average = os.path.join(path, "coverage.pdf")
    axis_plot(file_figure_average, x_axis, y_axis)
    file_figure_all = os.path.join(path, "all.pdf")
    axises_plot(file_figure_all, x_axis, y_axises, labels)

    uncovered_address = []
    for d in data:
        for u in d.uncovered_address:
            if not u in uncovered_address:
                uncovered_address.append(u)

    for d in data:
        for u in uncovered_address:
            if not u in d.uncovered_address:
                uncovered_address.remove(u)

    file_result = os.path.join(path, name_data_result)
    f = open(file_result, "w")
    f.write(str(uncovered_address))
    f.close()
    return stats, x_axis, y_axis


def dev_deal(dir_path, dir_name):
    stats = []
    x_axises = []
    y_axises = []
    labels = []

    path_dev = os.path.join(dir_path, dir_name)

    path_with_dra = os.path.join(path_dev, name_with_dra)
    result_get_from_sub_dir(stats, x_axises, y_axises, labels, path_with_dra, name_with_dra)

    path_without_dra = os.path.join(path_dev, name_without_dra)
    result_get_from_sub_dir(stats, x_axises, y_axises, labels, path_without_dra, name_without_dra)

    x_axises, y_axises = axis_expansion(length, x_axises, y_axises)
    x_axis = [sum(e) / len(e) for e in zip(*x_axises)]
    file_figure_all = os.path.join(dir_path, dir_name, dir_name + ".pdf")
    axises_plot(file_figure_all, x_axis, y_axises, labels)


def get_stat_file(path):
    is_dev = False
    is_results = False
    is_result = False

    dir_name = os.path.basename(path)
    dir_path = os.path.dirname(path)
    if dir_name.startswith(name_dev):
        dev_deal(dir_path, dir_name)
    elif dir_name.startswith(name_with_dra) or dir_name.startswith(name_without_dra):
        path_results = os.path.join(dir_path, dir_name)
        results_deal(path_results)
    elif dir_name.startswith(name_stat):
        result_deal(dir_path)
    else:
        for (dir_path, dir_names, file_names) in os.walk(path):
            for dir_name in dir_names:
                if dir_name.startswith(name_dev):
                    is_dev = True
                    dev_deal(dir_path, dir_name)
            if is_dev:
                break

            for dir_name in dir_names:
                if dir_name.startswith(name_with_dra) or dir_name.startswith(name_without_dra):
                    is_dev = True
                    path_results = os.path.join(dir_path, dir_name)
                    results_deal(path_results)
            if is_results:
                break

            for file_name in file_names:
                if file_name.startswith(name_stat):
                    is_result = True
                    result_deal(dir_path)
            if is_result:
                break


if __name__ == "__main__":
    if len(sys.argv) > 2:
        do_figure = False
    get_stat_file(sys.argv[1])
