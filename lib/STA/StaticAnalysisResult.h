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


    class Mod;
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

        MODS *GetAllGlobalWriteInsts(llvm::BasicBlock *B, bool branch);

        MODS *GetAllGlobalWriteInsts(BR_INF *p_taint_inf, bool branch);

        MODS *GetAllGlobalWriteBBs(llvm::BasicBlock *B, bool branch);

        MODS *GetAllGlobalWriteBBs(BR_INF *p_taint_inf, bool branch);

        std::string &getBBStrID(llvm::BasicBlock *B);

        std::string &getInstStrID(llvm::Instruction* I);

        std::string &getValueStr(llvm::Value *v);

        std::string &getTypeStr(llvm::Type *);

        static void stripFuncNameSuffix(std::string *fn);

        static llvm::DILocation *getCorrectInstrLocation(llvm::Instruction *I);

        TRAIT *getTrait(ID_TY);

        //Calculate the Levenshtein distance between two strings as a measure of fuzzy matching.
        static int levDistance(const std::string& source, const std::string& target);

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

    //A BB/Inst that can modify a global state.
    class Mod {
    public:
        Mod() {}

        Mod(llvm::BasicBlock *b, MOD_INF *pm, StaticAnalysisResult *sta) {
            this->init_properties();
            this->B = b;
            this->I = nullptr;
            this->mod_inf = *pm;
            this->pallcmds = nullptr;
            this->sta = sta;
        }

        Mod(llvm::Instruction *i, MOD_INF *pm, StaticAnalysisResult *sta) {
            this->init_properties();
            this->I = i;
            this->B = nullptr;
            if (i) {
                this->B = i->getParent();
            }
            this->mod_inf = *pm;
            this->pallcmds = nullptr;
            this->sta = sta;
        }

        ~Mod() {
            //
        }

        bool equal(const Mod *m) {
            if (!m) {
                return false;
            }
            return (this->B == m->B && this->I == m->I);
        }

        int calcPrio(std::string& cond, int64_t v) {
            int p = 0;
            TRAIT *pt = this->getSingleTrait();
            if ((!pt) || pt->empty()) {
                p = 0;
            }else if (this->from_nlp) {
                //TODO: What priority should we set for the mod IR from callee name NLP analysis?
                p = 0;
            }else if (cond == "==") {
                p = calcPrio_E(v);
            }else if (cond == "!=") {
                p = calcPrio_NE(v);
            }else if (cond == ">=" || cond == ">") {
                p = calcPrio_B(v, cond == ">=");
            }else if (cond == "<=" || cond == "<") {
                p = calcPrio_S(v, cond == "<=");
            }
            this->prio = p;
            return p;
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

        void init_properties() {
            this->sta = nullptr;
            this->repeat = 1;
            this->prio = 0;
            this->from_nlp = false;
            this->single_trait_id = 0;
        }

        llvm::BasicBlock *B;
        llvm::Instruction *I;
        //0: repeat but not sure about the times, >0: repeat at least for a certain time ("1" means no repeat).
        int64_t repeat;
        int prio;
        bool from_nlp;

    private:
        StaticAnalysisResult *sta;
        MOD_INF mod_inf;
        std::set<uint64_t> *pallcmds;
        uint64_t single_trait_id;
        TRAIT single_trait;

        //TODO: now we assume all traits are the same even under differnt contexts.
        //So only return one trait id.
        ID_TY getSingleTraitID() {
            if (this->single_trait_id) {
                return this->single_trait_id;
            }
            if (this->mod_inf.empty()) {
                return 0;
            }
            for (auto& x : this->mod_inf) {
                if (x.second.find(TRAIT_INDEX) == x.second.end()) {
                    continue;
                }
                std::set<uint64_t> &tids = x.second[TRAIT_INDEX];
                if (tids.empty()) {
                    continue;
                }
                for (auto& y : tids) {
                    this->single_trait_id = y;
                    return y;
                }
            }
            return 0;
        }

        TRAIT *getSingleTrait() {
            if (!this->single_trait.empty()) {
                return &(this->single_trait);
            }
            ID_TY stid = this->getSingleTraitID();
            if ((!this->sta) || (!stid)) {
                return nullptr;
            }
            TRAIT *pt = this->sta->getTrait(stid);
            if (pt) {
                this->single_trait = *pt;
                return &(this->single_trait);
            }
            return pt;
        }

        int calcPrio_E(int64_t n) {
            int p = 0;
            for (auto& x : this->single_trait) {
                std::string s = x.first;
                int64_t v = x.second;
                if (s == "CONST_INT") {
                    if (v == n) {
                        //This mod can (potentially) set the global state to the target value.
                        p = 100;
                    }else {
                        //Trait analysis is just a simple (and maybe inaccurate) pattern matching, so even the destination value is different,
                        //it's still possible to successfully set the global state, just give it a normal priority.
                        p = 0;
                    }
                }else {
                    //The modification is possibly accumulative, though we are not sure whether it will eventually set the global state as desired.
                    p = 0;
                    this->repeat = 0;
                }
                //Assume only one key exists in the trait.
                break;
            }
            return p;
        }

        int calcPrio_NE(int64_t n) {
            int p = 0;
            for (auto& x : this->single_trait) {
                std::string s = x.first;
                int64_t v = x.second;
                if (s == "CONST_INT") {
                    if (v != n) {
                        //This mod can (potentially) set the global state to a different target value.
                        p = 100;
                    }else {
                        //Trait analysis is just a simple (and maybe inaccurate) pattern matching, so even the destination value is the same,
                        //it's still possible to successfully set the global state as needed, just give it a normal priority.
                        p = 0;
                    }
                }else {
                    //The modification is possibly accumulative, though we are not sure whether it will eventually set the global state as desired.
                    //Since the condition is "!=", such a mod IR should be able to change the global state and satisfy the condition.
                    if (s == "ADD" || s == "SUB" || s == "MUL" || s == "DIV") {
                        p = 100;
                    }else {
                        p = 50;
                    }
                    this->repeat = 0;
                }
                //Assume only one key exists in the trait.
                break;
            }
            return p;
        }

        int calcPrio_B(int64_t n, bool inclusive) {
            int p = 0;
            for (auto& x : this->single_trait) {
                std::string s = x.first;
                int64_t v = x.second;
                if (s == "CONST_INT") {
                    if (v > n) {
                        //This mod can (potentially) set the global state to a bigger target value.
                        p = 100;
                    }else if (inclusive && v == n) {
                        p = 100;
                    }else {
                        //Trait analysis is just a simple (and maybe inaccurate) pattern matching, so even the destination value is the same,
                        //it's still possible to successfully set the global state as needed, just give it a normal priority.
                        p = 0;
                    }
                }else {
                    //The modification is possibly accumulative, though we are not sure whether it will eventually set the global state as desired.
                    //Since the condition is ">/>=", we need to exclude those mod IRs that will decrease the global states.
                    if (s == "ADD") {
                        p = (v > 0 ? 100 : -100);
                        //We are not sure how many times to repeat, though, since we don't know current value of the global state.
                        this->repeat = 0;
                    }else if (s == "SUB") {
                        p = (v < 0 ? 100 : -100);
                        //We are not sure how many times to repeat, though, since we don't know current value of the global state.
                        this->repeat = 0;
                    }else {
                        p = 0;
                        this->repeat = 0;
                    }
                }
                //Assume only one key exists in the trait.
                break;
            }
            return p;
        }

        int calcPrio_S(int64_t n, bool inclusive) {
            int p = 0;
            for (auto& x : this->single_trait) {
                std::string s = x.first;
                int64_t v = x.second;
                if (s == "CONST_INT") {
                    if (v < n) {
                        //This mod can (potentially) set the global state to a smaller target value.
                        p = 100;
                    }else if (inclusive && v == n) {
                        p = 100;
                    }else {
                        //Trait analysis is just a simple (and maybe inaccurate) pattern matching, so even the destination value is the same,
                        //it's still possible to successfully set the global state as needed, just give it a normal priority.
                        p = 0;
                    }
                }else {
                    //The modification is possibly accumulative, though we are not sure whether it will eventually set the global state as desired.
                    //Since the condition is ">/>=", we need to exclude those mod IRs that will decrease the global states.
                    if (s == "ADD") {
                        p = (v < 0 ? 100 : -100);
                        //We are not sure how many times to repeat, though, since we don't know current value of the global state.
                        this->repeat = 0;
                    }else if (s == "SUB") {
                        p = (v > 0 ? 100 : -100);
                        //We are not sure how many times to repeat, though, since we don't know current value of the global state.
                        this->repeat = 0;
                    }else {
                        p = 0;
                        this->repeat = 0;
                    }
                }
                //Assume only one key exists in the trait.
                break;
            }
            return p;
        }

    };

} /* namespace sta */

#endif /* LIB_STA_STATICANALYSISRESULT_H_ */
