/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#ifndef LIB_STA_STATICANALYSISRESULT_H_
#define LIB_STA_STATICANALYSISRESULT_H_

#include <llvm/IR/Module.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/Instruction.h>
#include <fstream>
#include <set>
#include "../JSON/json.cpp"

typedef std::vector<std::string> LOC_INF;
//ctx_id -> arg_no -> value set
typedef std::map<unsigned long,std::map<unsigned, std::set<uint64_t>>> MOD_INF;
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

        LOC_INF *getLocInf(llvm::Instruction *);

        LOC_INF *getLocInf(llvm::BasicBlock *);

        llvm::Instruction *getInstFromStr(std::string mod, std::string func, std::string bb, std::string inst);

        llvm::BasicBlock *getBBFromStr(std::string mod, std::string func, std::string bb);

        llvm::Module *p_module;

        MOD_IRS *GetAllGlobalWriteInsts(llvm::BasicBlock* B);

        MOD_IRS *GetAllGlobalWriteInsts(nlohmann::json *pj_taint_inf);

        MOD_BBS *GetAllGlobalWriteBBs(llvm::BasicBlock* B);

        MOD_BBS *GetAllGlobalWriteBBs(nlohmann::json *pj_taint_inf);

        //This is a temporary function...
        std::set<uint64_t> *getIoctlCmdSet(MOD_INF*);

    private:
        nlohmann::json j_taintedBrs, j_analysisCtxMap, j_tagMap, j_modInstCtxMap;

        nlohmann::json *findLocInJson(LOC_INF *p_loc, unsigned int e, nlohmann::json *data);

        nlohmann::json *QueryBranchTaint(llvm::BasicBlock* B);

        nlohmann::json *QueryModIRsFromTagID(unsigned long tid);

        nlohmann::json *QueryModIRsFromTagTy(std::string ty);

        MOD_IRS *j2ModIrs(nlohmann::json *pj_mod_irs);

        MOD_BBS *j2ModBbs(nlohmann::json *pj_mod_irs);
    };

} /* namespace sta */

#endif /* LIB_STA_STATICANALYSISRESULT_H_ */

