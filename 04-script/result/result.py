#! /usr/bin/python3
import os
import statistics
import sys

import scipy.stats
from result import data
from result import stats
from result import default
from result import axis


class device:
    def __init__(self, dir_path, dir_name):
        self.dir_path = dir_path
        self.path_dev = os.path.join(self.dir_path, dir_name)

        self.path_with_dra = os.path.join(self.path_dev, default.name_with_dra)
        self.results_with_dra = results(self.path_with_dra, 'C0')
        self.results_with_dra.deal()
        self.path_without_dra = os.path.join(self.path_dev, default.name_without_dra)
        self.results_without_dra = results(self.path_without_dra, 'C1')
        self.results_without_dra.deal()

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

    def get_coverage(self):
        max_coverage_with_dra = {}
        for s in self.results_with_dra.statistics.statistics:
            for a in s.coverage.coverage:
                if a not in max_coverage_with_dra:
                    max_coverage_with_dra[a] = 0
                else:
                    max_coverage_with_dra[a] = max_coverage_with_dra[a] + 1
        max_coverage_without_dra = {}
        for s in self.results_without_dra.statistics.statistics:
            for a in s.coverage.coverage:
                if a not in max_coverage_without_dra:
                    max_coverage_without_dra[a] = 0
                else:
                    max_coverage_without_dra[a] = max_coverage_without_dra[a] + 1

        unique_coverage_with_dra = {}
        for a in max_coverage_with_dra:
            if a not in max_coverage_without_dra:
                unique_coverage_with_dra[a] = max_coverage_with_dra[a]

        unique_coverage_without_dra = {}
        for a in max_coverage_without_dra:
            if a not in max_coverage_with_dra:
                unique_coverage_without_dra[a] = max_coverage_without_dra[a]

        max_coverage = {}
        for a in max_coverage_without_dra:
            max_coverage[a] = max_coverage_without_dra[a]
        for a in max_coverage_with_dra:
            if a not in max_coverage:
                max_coverage[a] = 0
            else:
                max_coverage[a] = max_coverage[a] + max_coverage_with_dra[a]

        file_result = os.path.join(self.path_dev, default.name_data_result)
        f = open(file_result, "w")
        f.write("unique_coverage_with_dra : " + str(len(unique_coverage_with_dra)) + "\n")
        f.write(str(unique_coverage_with_dra))
        f.write("unique_coverage_without_dra : " + str(len(unique_coverage_without_dra)) + "\n")
        f.write(str(unique_coverage_without_dra))
        f.write("max_coverage : " + str(len(max_coverage)) + "\n")

        f.close()


class results:
    def __init__(self, dir_path, color=''):
        self.dir_path = dir_path
        self.color = color
        self.results = []
        self.statistics = stats.stats(self.dir_path)
        self.axises = axis.axises(self.dir_path, self.color)

    def deal(self):
        if os.path.exists(self.dir_path):
            for (dir_path, dir_names, file_names) in os.walk(self.dir_path):
                for file_name in file_names:
                    if file_name.startswith(default.name_stat):
                        r = result(dir_path)
                        self.results.append(r)
                        self.statistics.statistics.append(r.stat)
                        self.axises.axises.append(r.axis)

        file_result = os.path.join(self.dir_path, default.name_stat_result)
        f = open(file_result, "w")
        self.statistics.get_average()
        f.write(self.statistics.stat.stat)
        self.statistics.stat.deal()
        f.write(self.statistics.stat.stat_deal)
        f.close()


class result:
    def __init__(self, dir_path):
        self.dir_path = dir_path
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
        device(dir_path, dir_name)
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
                    device(dir_path, dir_name)
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
