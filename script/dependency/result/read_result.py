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
length = 1 * 24 * 60
time_run = length * 60  # second


def stat_read(dir_path, file_name):
    # Read the existing Statistics.
    file_stat = os.path.join(dir_path, file_name)
    stat = pb.Statistics()
    f = open(file_stat, "rb")
    stat.ParseFromString(f.read())
    f.close()
    file_result = os.path.join(dir_path, name_stat_result)
    f = open(file_result, "w")
    f.write(str(stat))
    f.write(str(stat_deal(stat)))
    f.close()

    return stat


def read_results(path):
    stats = []
    if os.path.exists(path):
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.startswith(name_stat):
                    path_dev = os.path.join(dir_path, file_name)
                    stats.append(stat_read(path_dev))

    return stats


def plot_result(name, x_axis, y_axis):
    if not do_figure:
        return
    f = plt.figure()
    plt.plot(x_axis, y_axis)
    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    f.savefig(fname=name, bbox_inches='tight', format="pdf")


def plot_results(name, x_axis, y_axises, labels):
    if not do_figure:
        return
    f = plt.figure()
    for i in range(len(labels)):
        plt.plot(x_axis, y_axises[i], label=labels[i])

    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    plt.legend()
    f.savefig(fname=name, bbox_inches='tight', format="pdf")


def stat_get_time_coverage(stat):
    x_axis = []
    y_axis = []
    t0 = 0
    for i in stat.coverage.time:
        if i.time > t0:
            t0 = t0 + 60
            x_axis.append(t0)
            y_axis.append(i.num)
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


def get_average(stats):
    r = pb.Statistics()
    for stat in stats:
        for s in stat.stat:
            r.stat[s].name = stat.stat[s].name
            s_add(r.stat[s], stat.stat[s])
    length = len(stats)
    for s in r.stat:
        r.stat[s].executeNum = int(r.stat[s].executeNum/length)
        r.stat[s].time = int(r.stat[s].time/length)
        r.stat[s].newTestCaseNum = int(r.stat[s].newTestCaseNum/length)
        r.stat[s].newAddressNum = int(r.stat[s].newAddressNum/length)
    return r


def result_deal(dir_path, file_name):
    stat = stat_read(dir_path, file_name)
    x_axis, y_axis = stat_get_time_coverage(stat)
    file_figure = os.path.join(dir_path, "coverage.pdf")
    plot_result(file_figure, x_axis, y_axis)
    return stat, x_axis, y_axis


def expansion_axis(length, x_axises, y_axises):
    for x in x_axises:
        if len(x) != 0:
            max_time = x[-1]
        else:
            max_time = 0
        for i in range(length - len(x)):
            max_time = max_time + 60
            x.append(max_time)

    for y in y_axises:
        max_num = y[-1]
        for i in range(length - len(y)):
            y.append(max_num)
    return x_axises, y_axises


def deal_results(path):
    stats = []
    x_axises = []
    y_axises = []
    labels = []
    if os.path.exists(path):
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.startswith(name_stat):
                    s, x, y = result_deal(dir_path, file_name)
                    stats.append(s)
                    x_axises.append(x)
                    y_axises.append(y)
                    labels.append(dir_path)

    stat = get_average(stats)
    file_result = os.path.join(path, name_stat_result)
    f = open(file_result, "w")
    f.write(str(stat))
    f.write(str(stat_deal(stat)))
    f.close()

    x_axises, y_axises = expansion_axis(length, x_axises, y_axises)
    x_axis = [sum(e) / len(e) for e in zip(*x_axises)]
    y_axis = [sum(e) / len(e) for e in zip(*y_axises)]
    file_figure_average = os.path.join(path, "coverage.pdf")
    plot_result(file_figure_average, x_axis, y_axis)
    file_figure_all = os.path.join(path, "all.pdf")
    plot_results(file_figure_all, x_axis, y_axises, labels)
    return stats, x_axis, y_axis


def deal_dev(dir_path, dir_name):
    stats = []
    x_axises = []
    y_axises = []
    labels = []

    path_with_dra = os.path.join(dir_path, dir_name, name_with_dra)
    if os.path.exists(path_with_dra):
        s, x, y = deal_results(path_with_dra)
        stats.append(s)
        x_axises.append(x)
        y_axises.append(y)
        labels.append(name_with_dra)

    path_without_dra = os.path.join(dir_path, dir_name, name_without_dra)
    if os.path.exists(path_without_dra):
        s, x, y = deal_results(path_without_dra)
        stats.append(s)
        x_axises.append(x)
        y_axises.append(y)
        labels.append(name_without_dra)

    x_axises, y_axises = expansion_axis(length, x_axises, y_axises)
    x_axis = [sum(e) / len(e) for e in zip(*x_axises)]
    file_figure_all = os.path.join(dir_path, dir_name, dir_name + ".pdf")
    plot_results(file_figure_all, x_axis, y_axises, labels)


def get_stat_file(path):
    is_dev = False
    is_results = False
    is_result = False

    for (dir_path, dir_names, file_names) in os.walk(path):
        for dir_name in dir_names:
            if dir_name.startswith(name_dev):
                is_dev = True
                deal_dev(dir_path, dir_name)
        if is_dev:
            break

        for dir_name in dir_names:
            if dir_name.startswith(name_with_dra) or dir_name.startswith(name_without_dra):
                is_dev = True
                path_results = os.path.join(dir_path, dir_name)
                deal_results(path_results)
        if is_results:
            break

        for file_name in file_names:
            if file_name.startswith(name_stat):
                is_result = True
                result_deal(dir_path, file_name)
        if is_result:
            break


if __name__ == "__main__":
    if len(sys.argv) > 2:
        do_figure = False
    get_stat_file(sys.argv[1])
