/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#ifndef LIB_STA_STATICANALYSISRESULT_H_
#define LIB_STA_STATICANALYSISRESULT_H_

#include "../DRA/DBasicBlock.h"
#include <llvm/IR/Module.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/Instruction.h>
#include <fstream>
#include <set>
#include "../JSON/json.cpp"
#include "ResType.h"

typedef std::map<llvm::Instruction*,MOD_INF> MOD_IRS;
typedef std::map<llvm::BasicBlock*,MOD_INF> MOD_BBS;

namespace sta {

    class StaticAnalysisResult {
    public:
        StaticAnalysisResult(const std::string &staticRes, llvm::Module *p_module) {
            this->initStaticRes(staticRes, p_module);
        }

        StaticAnalysisResult() {
            //
        }

        virtual ~StaticAnalysisResult();

        int initStaticRes(const std::string &staticRes, llvm::Module *p_module);

        LOC_INF *getLocInf(dra::DLInstruction *);

        LOC_INF *getLocInf(dra::DBasicBlock *);

        llvm::Instruction *getInstFromStr(std::string mod, std::string func, std::string bb, std::string inst);

        llvm::BasicBlock *getBBFromStr(std::string mod, std::string func, std::string bb);

        llvm::Module *p_module;

        MOD_IRS *GetAllGlobalWriteInsts(dra::DBasicBlock *B);

        MOD_IRS *GetAllGlobalWriteInsts(ACTX_TAG_MAP *p_taint_inf);

        MOD_BBS *GetAllGlobalWriteBBs(dra::DBasicBlock* B);

        MOD_BBS *GetAllGlobalWriteBBs(ACTX_TAG_MAP *p_taint_inf);

        std::string& getBBStrID(llvm::BasicBlock* B);

        std::string& getValueStr(llvm::Value *v);

        std::string& getTypeStr(llvm::Type*);

        //This is a temporary function...
        std::set<uint64_t> *getIoctlCmdSet(MOD_INF*);

    private:
        nlohmann::json j_taintedBrs, j_analysisCtxMap, j_tagModMap, j_tagInfo, j_modInstCtxMap;

        TAINTED_BR_TY taintedBrs;
        ANALYSIS_CTX_MAP_TY analysisCtxMap;
        TAG_MOD_MAP_TY tagModMap;
        TAG_INFO_TY tagInfo;
        MOD_INST_CTX_MAP_TY modInstCtxMap;

        ACTX_TAG_MAP *QueryBranchTaint(dra::DBasicBlock* B);

        void QueryModIRsFromTagTy(std::string ty);

        MOD_IRS *GetRealModIrs(MOD_IR_TY *p_mod_irs);

        MOD_BBS *GetRealModBbs(MOD_IR_TY *p_mod_irs);
    };

} /* namespace sta */

#endif /* LIB_STA_STATICANALYSISRESULT_H_ */
