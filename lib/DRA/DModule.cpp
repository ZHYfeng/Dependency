#include "DModule.h"

#include <llvm/IR/DebugInfoMetadata.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/Module.h>
#include <llvm/IRReader/IRReader.h>
#include <llvm/Support/SourceMgr.h>
#include <fstream>
#include <iostream>
#include <sstream>

#include "../RPC/a2l.pb.h"

#define PATH_SIZE 10000

namespace dra {

    DModule::DModule() : addr2line(new dra::address()) {
        Function.reserve(PATH_SIZE);

    }

    DModule::~DModule() = default;

    void DModule::ReadBC(std::string InputFilename) {
#if DEBUG
        std::cout << "*************************************************" << std::endl;
        std::cout << "****************ReadIR***************************" << std::endl;
#endif
        llvm::LLVMContext *cxts;
        llvm::SMDiagnostic Err;
        cxts = new llvm::LLVMContext[1];
        module = llvm::parseIRFile(InputFilename, Err, cxts[0]);

        if (!module) {
            std::cerr << "load module: " << InputFilename << " failed\n";
            exit(0);
        } else {
#if DEBUGBC
            std::cerr << "size : " << module->getNamedMDList().size() << "\n";
            for (auto &i : module->getNamedMDList()) {
                i.dump();
            }
#endif
            BuildLLVMFunction(module.get());
        }
    }

    void DModule::BuildLLVMFunction(llvm::Module *Module) {
        DFunction *function;
        for (auto &it : *Module) {
            llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
            it.getAllMetadata(MDs);
            for (auto &MD : MDs) {
#if DEBUGBC
                MD.second->dump();
#endif
                if (llvm::MDNode *N = MD.second) {
                    if (auto *SP = llvm::dyn_cast<llvm::DISubprogram>(N)) {
                        std::string Path = SP->getFilename().str();
                        std::string Line = std::to_string(SP->getLine());
                        std::string name = it.getName().str();
                        std::string FunctionName;
                        if (name.find('.') < name.size()) {
                            FunctionName = name.substr(0, name.find('.'));
                        } else {
                            FunctionName = name;
                        }
                        function = CheckRepeatFunction(Path, FunctionName, dra::FunctionKind::IR);

                        function->IRName = name;
                        function->InitIRFunction(&it);
                        function->parent = this;
                    }
                }
            }

        }
    }

    std::string dra::DModule::exec(std::string cmd) {
        std::string data;
        FILE *stream;
        const int max_buffer = 256;
        char buffer[max_buffer];
        cmd.append(" 2>&1");

        stream = popen(cmd.c_str(), "r");
        if (stream) {
            while (!feof(stream))
                if (fgets(buffer, max_buffer, stream) != nullptr)
                    data.append(buffer);
            pclose(stream);
        }
        return data;
    }

    void dra::DModule::ReadObjdump(std::string objdump) {
        std::string Line;
        std::string Addr;
        std::string FunctionName;
        std::string Path;
        std::string Cmd;
        std::string Inst;
        std::string Result;
        std::stringstream ss;
        unsigned int LineNum;
        unsigned int InsNum;
        unsigned long end;
        unsigned long start;

        DFunction *function;
#if DEBUG
        std::cout << "*************************************************" << std::endl;
        std::cout << "****************ReadObjdump**********************" << std::endl;
#endif
#if DEBUGOBJDUMP
        std::cout << "objdump :" << objdump << std::endl;
#endif

        std::string obj = objdump.substr(0, objdump.find(".objdump"));
        std::string output_file = obj + ".bin";
        std::fstream input(output_file, std::ios::in | std::ios::binary);
        if (!this->addr2line->ParseFromIstream(&input)) {
            std::cerr << "Failed to parse addr2line." << std::endl;
        }
        input.close();

        std::ifstream objdumpFile(objdump);
        InsNum = 0;
        LineNum = 0;
        FunctionName = "";
        if (objdumpFile.is_open()) {
            while (getline(objdumpFile, Line)) {
                LineNum++;
                if (!Line.empty()) {
#if DEBUGOBJDUMP
                    std::cout << "Line :" << Line << std::endl;
#endif
                    if (Line.find(">:") < Line.size()) {
                        //deal with function
#if DEBUGOBJDUMP
                        std::cout << ">: :" << std::endl;
#endif

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
                        // get path
                        if (this->addr2line->mutable_addr()->find(Addr) != this->addr2line->mutable_addr()->end()) {
                            Path = (*this->addr2line->mutable_addr())[Addr];
                        } else {
                            std::cerr << "Failed to get addr2line." << std::endl;
                        }
#if DEBUGOBJDUMP
                        std::cout << "o Path :" << Path << std::endl;
#endif
                        function = CheckRepeatFunction(Path, FunctionName, dra::FunctionKind::O);
                        function->Address = Addr;

                    } else {
                        //asm instruction
                        if (Line.at(0) == '.') {
#if DEBUGOBJDUMP
                            std::cout << "dot :" << std::endl;
#endif
                        } else if (Line.at(0) == 'D') {
#if DEBUGOBJDUMP
                            std::cout << "D :" << std::endl;
#endif
                        } else if (Line.size() - Line.find(':') <= 23) {
                            // deal with no asm
                        } else {
#if DEBUGOBJDUMP
                            std::cout << "inst :" << std::endl;
#endif
                            InsNum++;

                            ss.str("");
                            for (char i : Line) {
                                ss << i;
                            }
                            std::string TempLine = ss.str();

                            unsigned int TempStart;
                            for (TempStart = 0; TempLine.at(TempStart) == ' '; TempStart++) {

                            }
                            Addr = "";
                            for (char i = TempLine.at(TempStart); i != ':'; TempStart++, i = TempLine.at(TempStart)) {
                                Addr += i;
                            }

                            Inst = TempLine.substr(TempLine.find(':') + 24, TempLine.size());
#if DEBUGOBJDUMP
                            std::cout << "o Addr :" << Addr << std::endl;
                            std::cout << "o Inst :" << Inst << std::endl;
#endif

                            if (Inst.at(0) == 'c' && Inst.find("call") < Inst.size()) {

                                auto *inst = new DAInstruction();
                                inst->setAddr(Addr);
                                inst->OInst = Inst;
                                if (Inst.find("__sanitizer_cov_trace_pc") < Inst.size()) {

                                }

                                function->InstASM.push_back(inst);
                                function->CallInstNum++;
                            }
                            if (Inst.at(0) == 'j') {
                                function->JumpInstNum++;
                            }
                        }
                    }

                } else if (InsNum > 0) {
                    // need add a space line at the end of objdump file.
                    if (!FunctionName.empty()) {
#if DEBUGOBJDUMP
                        std::cout << "Line :" << std::endl;
                        std::cout << "FunctionName :" << FunctionName << std::endl;
                        std::cout << "InsNum :" << InsNum << std::endl;
#endif
                        function->InstNum = InsNum;
                        InsNum = 0;
                        FunctionName = "";
                    }
                }
            }
            objdumpFile.close();
        } else {
            std::cerr << "Unable to open objdump file " << objdump << "\n";
        }
    }

    void dra::DModule::ReadAsmSourceCode(std::string AssemblySourceCode) {
        std::string line;
        std::string Path;
        std::string FunctionName;
        std::string BasicBlockName;
        std::string Inst;
        std::stringstream ss;
        unsigned int LineNum;
        unsigned int InstNum = 0;
        unsigned int CallInstNum = 0;
        unsigned int COVNum;

        DFunction *function;
#if DEBUG
        std::cout << "*************************************************" << std::endl;
        std::cout << "****************ReadAsmSourceCode****************" << std::endl;
#endif

#if DEBUGASM
        std::cout << "AssemblySourceCode :" << AssemblySourceCode << std::endl;
#endif

        std::ifstream AssemblySourceCodeFile(AssemblySourceCode);
        LineNum = 0;
        COVNum = 0;
        if (AssemblySourceCodeFile.is_open()) {
            while (getline(AssemblySourceCodeFile, line)) {
                LineNum++;
                if (!line.empty()) {
#if DEBUGASM
                    std::cout << "line :" << line << std::endl;
#endif
                    switch (line.at(0)) {
                        case '.': {
                            //label
#if DEBUGASM
                            std::cout << "dot :" << std::endl;
#endif
                            if (line.find(".Lfunc_end") < line.size()) {
                                function->InstNum = InstNum;
#if DEBUG
                                if (CallInstNum != function->InstASM.size()) {
                                    std::cout << "--------------------------------------------" << std::endl;
                                    std::cout << "different function : " << std::endl;
                                    std::cout << "Path :" << Path << std::endl;
                                    std::cout << "FunctionName :" << FunctionName << std::endl;
                                    std::cout << "InstASM.size() :" << function->InstASM.size() << std::endl;
                                    std::cout << "CallInstNum :" << CallInstNum << std::endl;
                                    std::cout << "COVNum :" << COVNum << std::endl;
                                    for (auto i : function->InstASM) {
                                        std::cout << "OInst :" << i->OInst << std::endl;
                                        std::cout << "SInst :" << i->SInst << std::endl;
                                    }
                                } else {
                                    std::cout << "--------------------------------------------" << std::endl;
                                    std::cout << "same function : " << std::endl;
                                    std::cout << "Path :" << Path << std::endl;
                                    std::cout << "FunctionName :" << FunctionName << std::endl;
                                    std::cout << "InstASM.size() :" << function->InstASM.size() << std::endl;
                                    std::cout << "CallInstNum :" << CallInstNum << std::endl;
                                    std::cout << "COVNum :" << COVNum << std::endl;
                                    for (auto i : function->InstASM) {
                                        std::cout << "OInst :" << i->OInst << std::endl;
                                        std::cout << "SInst :" << i->SInst << std::endl;
                                    }
                                }
#endif
#if DEBUGASM
                                std::cout << "FunctionName :" << FunctionName << std::endl;
                                std::cout << "InstNum :" << InstNum << std::endl;
#endif
                                InstNum = 0;
                                CallInstNum = 0;
                                COVNum = 0;
                            } else if (line.find("# %") < line.size()) {

                                ss.str("");
                                for (unsigned long i = line.find('%') + 1; i < line.size(); i++) {
                                    ss << line.at(i);
                                }
                                BasicBlockName = ss.str();
                                if (function->BasicBlock.find(BasicBlockName) != function->BasicBlock.end()) {
                                } else {
                                    function->BasicBlock[BasicBlockName] = new DBasicBlock();
                                    (function->BasicBlock[BasicBlockName])->name = BasicBlockName;
                                }

                                (function->BasicBlock[BasicBlockName])->setAsmSourceCode(true);

#if DEBUGASM
                                std::cout << ". bb name :" << ss.str() << std::endl;
#endif
                            }
                            break;
                        }
                        case '#': {
                            // bb
#if DEBUGASM
                            std::cout << "sharp :" << std::endl;
#endif
                            if (line.find("# %") < line.size()) {

                                ss.str("");
                                for (unsigned long i = line.find('%') + 1; i < line.size(); i++) {
                                    if (line.at(i) == '%') {
                                        for (i++; i < line.size(); i++) {
                                            ss << line.at(i);
                                        }
                                        BasicBlockName = ss.str();

                                        if (function->BasicBlock.find(BasicBlockName) != function->BasicBlock.end()) {
                                        } else {
                                            function->BasicBlock[BasicBlockName] = new DBasicBlock();
                                            (function->BasicBlock[BasicBlockName])->name = BasicBlockName;
                                        }
                                        (function->BasicBlock[BasicBlockName])->setAsmSourceCode(true);
                                        break;
                                    }
                                }

#if DEBUGASM
                                std::cout << "# bb name :" << ss.str() << std::endl;
#endif
                            }
                            break;
                        }
                        case '	': {
#if DEBUGASM
                            std::cout << "tab :" << std::endl;
                            std::cout << "line.size() :" << line.size() << std::endl;
#endif
                            if (line.size() == 1) {

                            } else if (line.at(1) == '.') {
                                //get path
                                if (Path.empty() && !FunctionName.empty() && line.find('#') < line.size()) {

                                    ss.str("");
                                    for (unsigned long i = line.find('#') + 2; i < line.find(':'); i++) {
                                        ss << line.at(i);
                                    }
                                    Path = ss.str();
#if DEBUGASM
                                    std::cout << "s Path :" << Path << std::endl;
#endif

                                    function = CheckRepeatFunction(Path, FunctionName, dra::FunctionKind::S);
                                }
                            } else if (line.at(1) == '#') {

                            } else if (line.at(1) >= 'a' && line.at(1) <= 'z') {
                                //asm instruction

                                if (0 && line.find("lock;") < line.size()) {

                                } else {
                                    ss.str("");
                                    for (unsigned long i = 1; i < line.size(); i++) {
                                        ss << line.at(i);
                                    }
                                    Inst = ss.str();
#if DEBUGASM
                                    std::cout << "s Inst :" << Inst << std::endl;
#endif
                                    if (CallInstNum >= function->InstASM.size()) {
                                    } else {
                                        if (Inst.at(0) == 'c' && Inst.find("call") <= Inst.size()) {
                                            auto *inst = function->InstASM.at(CallInstNum);
                                            inst->SInst = Inst;
                                            inst->BasicBlockName = BasicBlockName;
                                            inst->parent = function->BasicBlock[BasicBlockName];
                                            function->BasicBlock[BasicBlockName]->InstASM.push_back(inst);
                                            if (Inst.find("__sanitizer_cov_trace_pc") <= Inst.size()) {
                                                (function->BasicBlock[BasicBlockName])->COVNum++;
                                                (function->BasicBlock[BasicBlockName])->address = inst->address;
                                                COVNum++;
#if DEBUGASM
                                                std::cout << "o inst :" << inst->OInst << std::endl;
#endif
                                            }
                                            CallInstNum++;
                                        } else if (Inst.at(0) == 'c') {
                                            function->JumpInstNum++;
                                        }
                                    }
                                    InstNum++;
                                }

                            }

                            break;
                        }
                        case ' ': {
                            //comment
#if DEBUGASM
                            std::cout << "space :" << std::endl;
#endif
                            break;
                        }
                        default: {
                            if (line.find(':') < line.size()) {
                                if (line.find("# @") < line.size()) {
                                    ss.str("");
                                    for (unsigned long i = 0; line.at(i) != ':'; i++) {
                                        ss << line.at(i);
                                    }
                                    FunctionName = ss.str();
                                    Path = "";
#if DEBUGASM
                                    std::cout << "FunctionName :" << FunctionName << std::endl;
#endif
                                }
                            }
                        }
                    }
                }
            }
            AssemblySourceCodeFile.close();
        } else {
            std::cerr << "Unable to open AssemblySourceCodeFile " << AssemblySourceCode << ">\n";
        }
#if DEBUG
        std::cout << "****************ReadAsmSourceCode****************" << std::endl;
#endif
    }

    void DModule::AddRepeatFunction(DFunction *function, FunctionKind kind) {

        if (function->isRepeat()) {

        } else {
            function->setRepeat(true);
            switch (kind) {
                case dra::FunctionKind::IR: {
                    RepeatBCFunction[function->IRName] = function;
                    break;
                }
                case dra::FunctionKind::O: {
                    RepeatOFunction[function->Address] = function;
                    break;
                }
                case dra::FunctionKind::S: {
                    RepeatSFunction[function->Path].insert(std::pair<std::string, DFunction *>(function->FunctionName, function));
                    //maybe they are same
                    break;
                }
                default: {
                    std::cerr << "error AddRepeatFunction" << ">\n";
                }
            }
        }

    }

    DFunction *DModule::CheckRepeatFunction(std::string Path, std::string FunctionName, FunctionKind kind) {
        DFunction *function;
        if ((Function.find(Path) != Function.end()) && (Function[Path].find(FunctionName) != Function[Path].end())) {
            function = Function[Path][FunctionName];
            switch (kind) {
                case dra::FunctionKind::IR: {
                    if (function->isIR()) {
                        AddRepeatFunction(function, kind);
#if DEBUG
                        std::cout << "ir repeat function : " << std::endl;
                        function->dump();
#endif
                        function = CreatFunction(Path, FunctionName, kind);
                        AddRepeatFunction(function, kind);
                    } else {
                        function->setKind(kind);
                    }
                    break;
                }
                case dra::FunctionKind::O: {
                    if (function->isObjudump()) {
                        AddRepeatFunction(function, kind);
#if DEBUG
                        std::cout << "o repeat function : " << std::endl;
                        function->dump();
#endif
                        function = CreatFunction(Path, FunctionName, kind);
                        AddRepeatFunction(function, kind);
                    } else {
                        function->setKind(kind);
                    }
                    break;
                }
                case dra::FunctionKind::S: {
                    if (function->isAsmSourceCode()) {
                        AddRepeatFunction(function, kind);
#if DEBUG
                        std::cout << "s repeat function : " << std::endl;
                        function->dump();
#endif
                        function = CreatFunction(Path, FunctionName, kind);
                        AddRepeatFunction(function, kind);
                    } else {
                        function->setKind(kind);
                    }
                    break;
                }
                default: {
                }
            }

        } else {
            function = CreatFunction(Path, FunctionName, kind);
        }

        return function;
    }

    DFunction *DModule::CreatFunction(std::string Path, std::string FunctionName, FunctionKind kind) {
        DFunction *function;
        function = new DFunction();
        Function[Path].insert(std::pair<std::string, DFunction *>(FunctionName, function));
        function->FunctionName = FunctionName;
        function->Path = Path;
        function->setKind(kind);
        return function;
    }

} /* namespace dra */
