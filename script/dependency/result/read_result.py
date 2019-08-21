#! /usr/bin/python3
import os
import sys

import matplotlib.pyplot as plt

sys.path.append(os.getcwd())
from script.dependency.dra import DependencyRPC_pb2 as pb

name_dev = "dev_"
name_with_dra = "result-with-dra"
name_without_dra = "result-without-dra"
name_stat = "statistics.bin"
length = 1 * 24 * 60
time_run = length * 60  # second

def read_stat(file_stat):
    # Read the existing Statistics.
    stat = pb.Statistics()
    f = open(file_stat, "rb")
    stat.ParseFromString(f.read())
    f.close()
    # print(stat.coverage)
    # for i in stat.stat:
    #     print(i)
    #     print(stat.stat[i])

    return stat


def read_results(path):
    stats = []
    if os.path.exists(path):
        for (dir_path, dir_names, file_names) in os.walk(path):
            for file_name in file_names:
                if file_name.startswith(name_stat):
                    path_dev = os.path.join(dir_path, file_name)
                    stats.append(read_stat(path_dev))

    return stats


def plot_result(name, x_axis, y_axis):
    f = plt.figure()
    plt.plot(x_axis, y_axis)
    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    f.savefig(fname=name, bbox_inches='tight', format="pdf")


def plot_results(name, x_axis, y_axises, labels):
    f = plt.figure()
    for i in range(len(labels)):
        plt.plot(x_axis, y_axises[i], label=labels[i])

    plt.xlabel('time:second')
    plt.ylabel('coverage:address number')
    plt.title(name)
    plt.legend()
    f.savefig(fname=name, bbox_inches='tight', format="pdf")


def get_time_coverage(stat):
    x_axis = []
    y_axis = []
    t0 = 0
    for i in stat.coverage.time:
        if i.time > t0:
            t0 = t0 + 60
            x_axis.append(t0)
            y_axis.append(i.num)
    return x_axis, y_axis


def deal_result(dir_path, file_name):
    file = os.path.join(dir_path, file_name)
    stat = read_stat(file)
    x_axis, y_axis = get_time_coverage(stat)
    file_figure = os.path.join(dir_path, "coverage.pdf")
    plot_result(file_figure, x_axis, y_axis)
    return stat, x_axis, y_axis


def expansion_axis(length, x_axises, y_axises):
    for x in x_axises:
        max_time = x[len(x) - 1]
        for i in range(length - len(x)):
            max_time = max_time + 60
            x.append(max_time)

    for y in y_axises:
        max_num = y[len(y) - 1]
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
                    s, x, y = deal_result(dir_path, file_name)
                    stats.append(s)
                    x_axises.append(x)
                    y_axises.append(y)
                    labels.append(dir_path)

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
                deal_result(dir_path, file_name)
        if is_result:
            break


if __name__ == "__main__":
    print(sys.path)
    get_stat_file(sys.argv[1])
