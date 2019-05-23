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
#include <llvm/IR/DebugInfoMetadata.h>
#include <llvm/IR/DebugInfo.h>
#include <llvm/IR/CFG.h>
#include <fstream>
#include <set>
#include "../JSON/json.cpp"
#include "ResType.h"
#include "../DRA/DataManagement.h"

//typedef std::map<llvm::Instruction *, MOD_INF> MOD_IRS;
//typedef std::map<llvm::BasicBlock *, MOD_INF> MOD_BBS;

namespace sta {

    //A BB/Inst that can modify a global state.
    class Mod {
    public:
        Mod() {}

        Mod(llvm::BasicBlock *b, MOD_INF *pm) {
            this->B = b;
            this->I = nullptr;
            this->repeat = 0;
            this->mod_inf = *pm;
            this->pallcmds = nullptr;
        }

        Mod(llvm::Instruction *i, MOD_INF *pm) {
            this->I = i;
            this->B = nullptr;
            if (i) {
                this->B = i->getParent();
            }
            this->repeat = 0;
            this->mod_inf = *pm;
            this->pallcmds = nullptr;
        }

        ~Mod() {
            //
        }

        std::set<uint64_t> *getIoctlCmdSet() {
            if (this->pallcmds) {
                return this->pallcmds;
            }
            if (this->mod_inf.empty()) {
                return nullptr;
            }
            this->pallcmds = new std::set<uint64_t>();
            for (auto &x : this->mod_inf) {
                std::set<uint64_t> &cs = x.second[1];
                this->pallcmds->insert(cs.begin(), cs.end());
            }
            return this->pallcmds;
        }

        llvm::BasicBlock *B;
        llvm::Instruction *I;
        int64_t repeat;

    private:
        MOD_INF mod_inf;
        std::set<uint64_t> *pallcmds;
    };

    typedef std::vector<Mod*> MODS;

    class StaticAnalysisResult {
    public:
        StaticAnalysisResult(const std::string &staticRes, dra::DataManagement *DM) {
            this->initStaticRes(staticRes, DM);
        }

        StaticAnalysisResult() {
            //
        }

        virtual ~StaticAnalysisResult();

        int initStaticRes(const std::string &staticRes, dra::DataManagement *DM);

        LOC_INF *getLocInf(llvm::Instruction *, bool);

        LOC_INF *getLocInf(llvm::BasicBlock *);

        llvm::Instruction *getInstFromStr(std::string path, std::string func, std::string bb, std::string inst);

        llvm::BasicBlock *getBBFromStr(std::string path, std::string func, std::string bb);

        llvm::Module *p_module;

        dra::DataManagement *dm;

        MODS *GetAllGlobalWriteInsts(llvm::BasicBlock *B);

        MODS *GetAllGlobalWriteInsts(BR_INF *p_taint_inf);

        MODS *GetAllGlobalWriteBBs(llvm::BasicBlock *B);

        MODS *GetAllGlobalWriteBBs(BR_INF *p_taint_inf);

        std::string &getBBStrID(llvm::BasicBlock *B);

        std::string &getInstStrID(llvm::Instruction* I);

        std::string &getValueStr(llvm::Value *v);

        std::string &getTypeStr(llvm::Type *);

        static void stripFuncNameSuffix(std::string *fn);

        static llvm::DILocation *getCorrectInstrLocation(llvm::Instruction *I);

        //This is a temporary function...
        std::set<uint64_t> *getIoctlCmdSet(MOD_INF *);

    private:
        nlohmann::json j_taintedBrs, j_ctxMap, j_traitMap, j_tagModMap, j_tagInfo, j_calleeMap;

        TAINTED_BR_TY taintedBrs;
        CTX_MAP_TY ctxMap;
        INST_TRAIT_MAP traitMap;
        TAG_MOD_MAP_TY tagModMap;
        TAG_INFO_TY tagInfo;
        CALLEE_MAP_TY calleeMap;

        BR_INF *QueryBranchTaint(llvm::BasicBlock *B);

        void QueryModIRsFromTagTy(std::string ty);

        MODS *GetRealModIrs(MOD_IR_TY *p_mod_irs);

        MODS *GetRealModBbs(MOD_IR_TY *p_mod_irs);

        void tweakModsOnTraits(MODS *pmods, ID_TY br_trait_id, bool branch);
    };

} /* namespace sta */

#endif /* LIB_STA_STATICANALYSISRESULT_H_ */
