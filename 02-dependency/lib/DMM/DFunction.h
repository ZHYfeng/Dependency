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

    static DFunction *MargeDFunction(DFunction *one, DFunction *two);

    class DFunction {
    public:
        DFunction();

        virtual ~DFunction();

        void InitIRFunction(llvm::Function *f);

        void setState(CoverKind kind);

        void update(CoverKind kind);

        bool isObjudump() const;

        void setObjudump(bool Objudump);

        bool isAsmSourceCode() const;

        void setAsmSourceCode(bool AsmSourceCode);

        bool isIR() const;

        void setIR(bool IR);

        void setKind(FunctionKind kind);

        bool isMap();

        bool isRepeat() const;

        void setRepeat(bool repeat);

        void dump();

        void inferUseLessPred(llvm::BasicBlock *b);

        void inferUseLessPred();

        void compute_arrive();

        void get_terminator(std::vector<dra::DBasicBlock *> &terminator_bb);

        void set_pred_successor(DBasicBlock *db);

        void set_critical_condition();

        uint32_t get_uncovered_basicblock_number();

        void get_function_call(std::set<llvm::Function *> &res);

    public:
        bool Objudump;
        bool AsmSourceCode;
        bool IR;

        bool repeat;

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
        unsigned int RealBasicBlockNum;
        unsigned int BasicBlockNum;
        std::unordered_map<std::string, DBasicBlock *> BasicBlock;

        std::vector<llvm::BasicBlock *> path;
        std::set<llvm::BasicBlock *> order;

        bool critical_condition;

        bool uncovered_basicblock;

    };

} /* namespace dra */

#endif /* LIB_DRA_FUNCTION_H_ */
