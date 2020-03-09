#! /usr/bin/python3
import os
import subprocess

import scipy.stats

import default
import read_axis
import read_stats


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
        self.axises = read_axis.axises(self.path_dev)
        self.set_axises()

        self.statistic, self.p_value = 0, 0
        self.get_mann_withney_utest()

        # self.get_coverage()
        # self.get_base()

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

class Results:
    def __init__(self, dir_path, color=''):
        self.dir_path = dir_path
        self.color = color

        self.file_result = os.path.join(self.dir_path, default.name_data_result)
        # if os.path.exists(self.file_result):
        #     os.remove(self.file_result)

        self.results = []
        self.statistics = read_stats.stats(self.dir_path)
        self.axises = read_axis.axises(self.dir_path, self.color)

        self.max_coverage = {}
        self.uncovered_address_input = []
        self.uncovered_address_dependency = []
        self.max_uncoverage = {}
        self.deal_results()
        # self.get_uncovered_address()
        # self.get_max_coverage()

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
        # print("self.stat = stats.stat(self.dir_path)")
        self.stat = read_stats.stat(self.dir_path)
        # print("self.stat.get_time_coverage()")
        self.stat.get_time_coverage()
        # print("self.axis = axis.axis(self.dir_path, self.stat.x_axis, self.stat.y_axis, '-')")
        self.axis = read_axis.axis(self.dir_path, self.stat.x_axis, self.stat.y_axis, '-')


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
