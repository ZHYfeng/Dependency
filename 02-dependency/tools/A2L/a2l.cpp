//
// Created by yhao on 5/15/19.
//

#include <llvm/Support/CommandLine.h>
#include <sstream>
#include <fstream>
#include "../../lib/RPC/Dependency.pb.h"

#define TEST 0
#define DEBUG_OBJ_DUMP 0

llvm::cl::opt<std::string> objdump("objdump", llvm::cl::desc("The path of objdump."), llvm::cl::init("./vmlinux.objdump"));

int main(int argc, char **argv) {

    llvm::cl::ParseCommandLineOptions(argc, argv, "A2L\n");
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

#if DEBUG_OBJ_DUMP
    std::cout << "objdump :" << objdump << std::endl;
#endif

    std::ifstream objdumpFile(objdump);
    FunctionName = "";
    if (objdumpFile.is_open()) {
        while (getline(objdumpFile, Line)) {
            if (!Line.empty()) {
                if (Line.find(">:") < Line.size()) {
                    //deal with function

                    // get trace_pc_address
                    ss.str("");
                    for (unsigned long i = 0; i < 16; i++) {
                        ss << Line.at(i);
                    }
                    Addr = ss.str();
#if DEBUG_OBJ_DUMP
                    std::cout << "o Addr :" << Addr << std::endl;
#endif


                    // get function name
                    ss.str("");
                    start = Line.find('<');
                    end = Line.find('>');
                    for (unsigned long i = start + 1; i < end; i++) {
                        ss << Line.at(i);
                    }
                    FunctionName = ss.str();
#if DEBUG_OBJ_DUMP
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
                                Result = buffer;
                        pclose(stream);
                    }
#if DEBUG_OBJ_DUMP
                    std::cout << "Result :" << Result << std::endl;
#endif
                    ss.str("");
#if TEST
                    start = Result.find("c-f/");
#else
                    start = Result.find("-np/");
//                    start = Result.find(".16/");
#endif
                    end = Result.find(':');
                    for (unsigned long i = start + 4; i < end; i++) {
                        ss << Result.at(i);
                    }
                    Path = ss.str();
#if DEBUG_OBJ_DUMP
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

    // Write back to disk.
    std::string output_file = obj + ".bin";
    std::fstream output(output_file, std::ios::out | std::ios::trunc | std::ios::binary);
    if (!a->SerializeToOstream(&output)) {
        std::cerr << "Failed to write map." << std::endl;
        return -1;
    }
    output.close();

    return 0;
}