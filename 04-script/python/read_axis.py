import os
import statistics

import matplotlib.pyplot as plt

from python import default


class axis:
    def __init__(self, dir_path, x_axis, y_axis, line_style):

        self.dir_path = dir_path
        self.x_axis = x_axis
        self.y_axis = y_axis
        self.line_style = line_style
        self.file_figure = os.path.join(self.dir_path, "coverage.pdf")
        self.expansion(default.length + 10)
        # self.plot()

    def plot(self):
        if not default.do_figure or len(self.y_axis) == 0:
            return
        f = plt.figure()
        plt.plot(self.x_axis, self.y_axis)
        plt.xlabel('time:second')
        plt.ylabel('coverage:address number')
        plt.title(self.file_figure)
        f.savefig(fname=self.file_figure, bbox_inches='tight', format="pdf")
        plt.close(f)

    def expansion(self, length):

        if len(self.x_axis) != 0:
            max_time = self.x_axis[-1]
        else:
            max_time = 0
        for i in range(length - len(self.x_axis)):
            max_time = max_time + 60
            self.x_axis.append(max_time)

        if len(self.y_axis) != 0:
            max_num = self.y_axis[-1]
        else:
            max_num = 0
        for i in range(length - len(self.y_axis)):
            self.y_axis.append(max_num)


class axises:
    def __init__(self, dir_path, color=''):
        self.dir_path = dir_path
        self.color = color

        self.axises = []
        self.x_axis = []
        self.x_axises = []

        self.y_axises = []
        self.labels = []
        self.line_styles = []
        self.colors = []

        self.y_axises_statistics = []
        self.labels_statistics = []
        self.line_styles_statistics = []
        self.colors_statistics = []

    def plot(self, name, y_axises=None, labels=None, line_styles=None, colors=None, title=""):
        if y_axises is None:
            y_axises = self.y_axises
        if labels is None:
            labels = self.labels
        if line_styles is None:
            line_styles = self.line_styles
        if colors is None:
            colors = self.colors

        if not default.do_figure or len(y_axises) == 0:
            return
        f = plt.figure()
        if len(colors) == 0:
            for i in range(len(labels)):
                plt.plot(self.x_axis, y_axises[i], label=labels[i], linestyle=line_styles[i])
        else:
            for i in range(len(labels)):
                plt.plot(self.x_axis, y_axises[i], label=labels[i], linestyle=line_styles[i], color=colors[i])

        plt.xlabel('time:second')
        plt.ylabel('coverage:address number')
        plt.title(name + title)
        if len(labels) != 0:
            plt.legend()
        f.savefig(fname=name, bbox_inches='tight', format="pdf")
        plt.close(f)

    def deal(self):

        for a in self.axises:
            self.x_axises.append(a.x_axis)
            self.y_axises.append(a.y_axis)
            self.line_styles.append(a.line_style)
            self.labels.append(a.dir_path)

        self.x_axis = [sum(e) / len(e) for e in zip(*self.x_axises)]

        y_axis_mean = [statistics.mean(e) for e in zip(*self.y_axises)]
        self.y_axises_statistics.append(y_axis_mean)
        self.labels_statistics.append("mean")
        self.line_styles_statistics.append(':')
        if not self.color == '':
            self.colors_statistics.append(self.color)

        # y_axis_median = [statistics.median(e) for e in zip(*self.y_axises)]
        # self.y_axises_statistics.append(y_axis_median)
        # self.labels_statistics.append("median")
        # self.line_styles_statistics.append('-')
        # if not self.color == '':
        #     self.colors_statistics.append(self.color)

        y_axis_max = [max(e) for e in zip(*self.y_axises)]
        self.y_axises_statistics.append(y_axis_max)
        self.labels_statistics.append("max")
        self.line_styles_statistics.append('--')
        if not self.color == '':
            self.colors_statistics.append(self.color)

        y_axis_min = [min(e) for e in zip(*self.y_axises)]
        self.y_axises_statistics.append(y_axis_min)
        self.labels_statistics.append("min")
        self.line_styles_statistics.append('--')
        if not self.color == '':
            self.colors_statistics.append(self.color)

        # y_axis_confidence_intervals_start = [
        #     scipy.mean(e) - scipy.stats.sem(e) * scipy.stats.t.ppf((1 + devices.confidence) / 2, len(e) - 1) for e in
        #     zip(*self.y_axises)]
        # self.y_axises_statistics.append(y_axis_confidence_intervals_start)
        # self.labels_statistics.append("ci_start")
        # self.line_styles_statistics.append('-.')
        # if not self.color == '':
        #     self.colors_statistics.append(self.color)
        #
        # y_axis_confidence_intervals_end = [
        #     scipy.mean(e) + scipy.stats.sem(e) * scipy.stats.t.ppf((1 + devices.confidence) / 2, len(e) - 1) for e in
        #     zip(*self.y_axises)]
        # self.y_axises_statistics.append(y_axis_confidence_intervals_end)
        # self.labels_statistics.append("ci_end")
        # self.line_styles_statistics.append('-.')
        # if not self.color == '':
        #     self.colors_statistics.append(self.color)

        if not len(self.y_axises) == 0:
            file_figure_average = os.path.join(self.dir_path, "coverage.pdf")
            self.plot(file_figure_average, self.y_axises_statistics, self.labels_statistics,
                      self.line_styles_statistics, self.colors_statistics)

            file_figure_all = os.path.join(self.dir_path, "all.pdf")
            self.plot(name=file_figure_all)
