/*
 * DFunction.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_FUNCTION_H_
#define LIB_DRA_FUNCTION_H_

#include <llvm/IR/BasicBlock.h>
#include <string>
#include <unordered_map>
#include <vector>

#include "DAInstruction.h"
#include "DBasicBlock.h"
#include "DLInstruction.h"

namespace dra {
    class DModule;
} /* namespace dra */

namespace llvm {
    class Function;
} /* namespace llvm */

namespace dra {

    enum FunctionKind {
        IR, O, S,
    };

    class DFunction {
    public:
        DFunction();

        virtual ~DFunction();

        void InitIRFunction(llvm::Function *f);

        void setState(CoverKind kind);

        void update(CoverKind kind);

        bool isObjudump() const;

        void setObjudump(bool objudump);

        bool isAsmSourceCode() const;

        void setAsmSourceCode(bool asmSourceCode);

        bool isIR() const;

        void setIR(bool ir);

        void setKind(FunctionKind kind);

        bool isMap() const;

        bool isRepeat() const;

        void setRepeat(bool repeat);

        void dump() const;

        uint32_t get_number_uncovered_instructions(std::set<dra::DBasicBlock *> &res);

        void get_function_call(std::set<llvm::Function *> &res);

        uint32_t get_uncovered_instructions(llvm::BasicBlock *b, std::set<dra::DBasicBlock *> &res);

        uint32_t get_dominator_uncovered_instructions(llvm::BasicBlock *b, std::set<dra::DBasicBlock *> &res);

        void add_number_basic_block_covered();

    public:
        bool Objudump;
        bool AsmSourceCode;
        bool IR;

        bool Repeat;


        llvm::Function *function;
        DModule *parent;
        CoverKind state;

        std::string FunctionName;
        std::string IRName;
        std::string Path;

        std::string Address;
        unsigned int InstNum;
        unsigned int CallInstNum;
        unsigned int JumpInstNum;
        std::vector<DAInstruction *> InstASM;
        unsigned int NumberBasicBlock;
        unsigned int NumberBasicBlockReal;
        std::unordered_map<std::string, DBasicBlock *> BasicBlock;

        std::vector<llvm::BasicBlock *> path;
        std::set<llvm::BasicBlock *> order;

        u_int32_t NumberBasicBlockCovered{};

        bool uncovered_basicblock;
        llvm::DominatorTree *DT{};
    };

} /* namespace dra */

#endif /* LIB_DRA_FUNCTION_H_ */
