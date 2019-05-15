//
// Created by yhao on 5/15/19.
//

#include <llvm/Support/CommandLine.h>
#include <sstream>
#include <fstream>
#include "../../lib/RPC/a2l.pb.cc"

#define TEST 1
#define DEBUGOBJDUMP 1

llvm::cl::opt<std::string> objdump("objdump", llvm::cl::desc("The path of objdump."), llvm::cl::init("./vmlinux.objdump"));

int main(int argc, char **argv) {

    llvm::cl::ParseCommandLineOptions(argc, argv, "a2l\n");
    std::unique_ptr<dra::address> a(new dra::address());

    std::string Line;
    std::string Addr;
    std::string FunctionName;
    std::string Path;
    std::string cmd;
    std::string Result;
    std::stringstream ss;
    unsigned long end;
    unsigned long start;

    // get path
    std::string obj = objdump.substr(0, objdump.find(".objdump"));

#if DEBUGOBJDUMP
    std::cout << "objdump :" << objdump << std::endl;
#endif

    std::ifstream objdumpFile(objdump);
    FunctionName = "";
    if (objdumpFile.is_open()) {
        while (getline(objdumpFile, Line)) {
            if (!Line.empty()) {
                if (Line.find(">:") < Line.size()) {
                    //deal with function

                    // get address
                    ss.str("");
                    for (unsigned long i = 0; i < 16; i++) {
                        ss << Line.at(i);
                    }
                    Addr = ss.str();

                    // get function name
                    ss.str("");
                    start = Line.find('<');
                    end = Line.find('>');
                    for (unsigned long i = start + 1; i < end; i++) {
                        ss << Line.at(i);
                    }
                    FunctionName = ss.str();
#if DEBUGOBJDUMP
                    std::cout << "o FunctionName :" << FunctionName << std::endl;
#endif

#if TEST
                    cmd = "addr2line -afi -e " + obj + ".o " + Addr;
#else
                    cmd = "addr2line -afi -e " + obj + " " + Addr;
#endif

                    FILE *stream;
                    const int max_buffer = 1024;
                    char buffer[max_buffer];
                    cmd.append(" 2>&1");

                    stream = popen(cmd.c_str(), "r");
                    if (stream) {
                        while (!feof(stream))
                            if (fgets(buffer, max_buffer, stream) != nullptr)
                                Result.append(buffer);
                        pclose(stream);
                    }

                    ss.str("");
                    start = Result.find("-np/");

                    end = Result.find(':');
                    for (unsigned long i = start + 4; i < end; i++) {
                        ss << Result.at(i);
                    }
                    Path = ss.str();
#if DEBUGOBJDUMP
                    std::cout << "o Path :" << Path << std::endl;
#endif
                    (*a->mutable_addr())[Addr] = Path;
                }
            }
        }
        objdumpFile.close();
    } else {
        std::cerr << "Unable to open objdump file " << objdump << "\n";
    }


    a->set_name("asd");
    std::cout << "a" << std::endl;
    for (auto c : *a->mutable_addr()){
        std::cout << c.first << " " << c.second << std::endl;
    }

    // Write back to disk.
    std::fstream output("./a2l.bin", std::ios::out | std::ios::trunc | std::ios::binary);
    if (!a->SerializeToOstream(&output)) {
        std::cerr << "Failed to write map." << std::endl;
        return -1;
    }

    // Read.
    std::unique_ptr<dra::address> b(new dra::address());
    std::fstream input("./a2l.bin", std::ios::in | std::ios::binary);
    if (!b->ParseFromIstream(&input)) {
        std::cerr << "Failed to parse map." << std::endl;
        return -1;
    }

    std::cout << "b" << b->name() << std::endl;
    for (auto c : b->addr()){
        std::cout << c.first << " " << c.second << std::endl;
    }

    return 0;
}