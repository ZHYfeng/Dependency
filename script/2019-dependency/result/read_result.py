import dra.DependencyRPC_pb2


def read_stat(file_name):
    # Read the existing Statistics.
    stat = dra.DependencyRPC_pb2.Statistics()
    f = open(file_name, "rb")
    stat.ParseFromString(f.read())
    f.close()
