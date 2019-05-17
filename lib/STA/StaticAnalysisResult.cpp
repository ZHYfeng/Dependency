/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#include "StaticAnalysisResult.h"
#include "../DRA/DFunction.h"
#include <iostream>

namespace sta {

    StaticAnalysisResult::~StaticAnalysisResult() = default;

    int StaticAnalysisResult::initStaticRes(const std::string &staticRes, llvm::Module *p_module) {
        this->p_module = p_module;
        try {
            std::ifstream infile;
            infile.open(staticRes);
            infile >> this->j_taintedBrs >> this->j_analysisCtxMap >> this->j_tagModMap >> this->j_tagInfo >> this->j_modInstCtxMap;
            infile.close();
            this->taintedBrs = this->j_taintedBrs.get<TAINTED_BR_TY>();
            this->analysisCtxMap = this->j_analysisCtxMap.get<ANALYSIS_CTX_MAP_TY>();
            this->tagModMap = this->j_tagModMap.get<TAG_MOD_MAP_TY>();
            this->tagInfo = this->j_tagInfo.get<TAG_INFO_TY>();
            this->modInstCtxMap = this->j_modInstCtxMap.get<MOD_INST_CTX_MAP_TY>();
            return 0;
        } catch (...) {
            std::cout << "Fail to deserialize the static analysis results!\n";
        }
        return 1;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(dra::DLInstruction *I) {
        if (!I) {
            return nullptr;
        }
        std::string inst, bb, func, file;
        std::string str;
        llvm::raw_string_ostream ss(str);
        ss << *I->i;
        inst = ss.str();

        auto B = I->parent;
        if (B != nullptr) {
            bb = B->name;
            auto F = B->parent;
            if(F != nullptr){
                func = F->FunctionName;
                file = F->Path;
            }
        }

        LOC_INF *loc_inf = new LOC_INF;
        loc_inf->push_back(inst);
        loc_inf->push_back(bb);
        loc_inf->push_back(func);
        loc_inf->push_back(file);
        return loc_inf;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(dra::DBasicBlock *B) {
        if (!B) {
            return nullptr;
        }
        return this->getLocInf(*(B->InstIR.begin()));
    }

    //Given a bb, return the taint information regarding its last br inst.
    //The returned info is a map from the context id to the taint tag id set.
    ACTX_TAG_MAP *StaticAnalysisResult::QueryBranchTaint(dra::DBasicBlock *B) {
        if (!B) {
            return nullptr;
        }
        LOC_INF *p_loc = this->getLocInf(B);
        if (!p_loc) {
            return nullptr;
        }
        auto &res3 = this->taintedBrs;
        if (res3.find((*p_loc)[3]) != res3.end()) {
            auto &res2 = res3[(*p_loc)[3]];
            if (res2.find((*p_loc)[2]) != res2.end()) {
                auto &res1 = res2[(*p_loc)[2]];
                if (res1.find((*p_loc)[1]) != res1.end()) {
                    return &(res1[(*p_loc)[1]]);
                }
            }
        }
        return nullptr;
    }

    MOD_IRS *StaticAnalysisResult::GetAllGlobalWriteInsts(dra::DBasicBlock *B) {
        return this->GetAllGlobalWriteInsts(this->QueryBranchTaint(B));
    }

    //Whatever call context under which the br is tainted, we will contain its mod insts for any tags (i.e. ALL).
    MOD_IRS *StaticAnalysisResult::GetAllGlobalWriteInsts(ACTX_TAG_MAP *p_taint_inf) {
        if (!p_taint_inf) {
            return nullptr;
        }
        MOD_IRS *p_mod_irs = new MOD_IRS();
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            auto &tag_ids = x.second;
            for (ID_TY tid : tag_ids) {
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MOD_IRS *p_cur_mod_irs = this->GetRealModIrs(ps_mod_irs);
                //Merge.
                for (auto const &x : *p_cur_mod_irs) {
                    if (p_mod_irs->find(x.first) != p_mod_irs->end()) {
                        (*p_mod_irs)[x.first].insert(x.second.begin(), x.second.end());
                    } else {
                        (*p_mod_irs)[x.first] = x.second;
                    }
                }//merge
            }//tags
        }
        return p_mod_irs;
    }

    MOD_BBS *StaticAnalysisResult::GetAllGlobalWriteBBs(dra::DBasicBlock *B) {
        return this->GetAllGlobalWriteBBs(this->QueryBranchTaint(B));
    }

    MOD_BBS *StaticAnalysisResult::GetAllGlobalWriteBBs(ACTX_TAG_MAP *p_taint_inf) {
        if (!p_taint_inf) {
            return nullptr;
        }
        MOD_BBS *p_mod_bbs = new MOD_BBS();
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            auto &tag_ids = x.second;
            for (ID_TY tid : tag_ids) {
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MOD_BBS *p_cur_mod_bbs = this->GetRealModBbs(ps_mod_irs);
                //Merge.
                for (auto const &x : *p_cur_mod_bbs) {
                    if (p_mod_bbs->find(x.first) != p_mod_bbs->end()) {
                        (*p_mod_bbs)[x.first].insert(x.second.begin(), x.second.end());
                    } else {
                        (*p_mod_bbs)[x.first] = x.second;
                    }
                }//merge
            }//tags
        }
        return p_mod_bbs;
    }

    MOD_IRS *StaticAnalysisResult::GetRealModIrs(MOD_IR_TY *p_mod_irs) {
        if (!p_mod_irs) {
            return nullptr;
        }
        MOD_IRS *mod_irs = new MOD_IRS();
        for (auto &el0 : *p_mod_irs) {
            const std::string &module = el0.first;
            for (auto &el1 : (*p_mod_irs)[module]) {
                const std::string &func = el1.first;
                for (auto &el2 : (*p_mod_irs)[module][func]) {
                    const std::string &bb = el2.first;
                    for (auto &el3 : (*p_mod_irs)[module][func][bb]) {
                        const std::string &inst = el3.first;
                        //Get the actual Instruction* according to these string info
                        llvm::Instruction *pinst = this->getInstFromStr(module, func, bb, inst);
                        if (!pinst) {
                            continue;
                        }
                        (*mod_irs)[pinst] = el3.second;
                    }//inst
                }//bb
            }//func
        }//module
        return mod_irs;
    }

    MOD_BBS *StaticAnalysisResult::GetRealModBbs(MOD_IR_TY *p_mod_irs) {
        if (!p_mod_irs) {
            return nullptr;
        }
        MOD_BBS *mod_bbs = new MOD_BBS();
        for (auto &el0 : *p_mod_irs) {
            const std::string &module = el0.first;
            for (auto &el1 : (*p_mod_irs)[module]) {
                const std::string &func = el1.first;
                for (auto &el2 : (*p_mod_irs)[module][func]) {
                    const std::string &bb = el2.first;
                    llvm::BasicBlock *pbb = this->getBBFromStr(module, func, bb);
                    if (!pbb) {
                        continue;
                    }
                    for (auto &el3 : (*p_mod_irs)[module][func][bb]) {
                        const MOD_INF &mod_inf = el3.second;
                        (*mod_bbs)[pbb].insert(mod_inf.begin(), mod_inf.end());
                    }//inst
                }//bb
            }//func
        }//module
        return mod_bbs;
    }

    llvm::Instruction *StaticAnalysisResult::getInstFromStr(std::string mod, std::string func, std::string bb, std::string inst) {
        //NOTE: Since now we only have one module, skip the module name match..
        for (llvm::Function &curFunc : *(this->p_module)) {
            if (curFunc.getName().str() != func) {
                continue;
            }
            for (llvm::BasicBlock &curBB : curFunc) {
                if (this->getBBStrID(&curBB) != bb) {
                    //if (curBB.getName().str() != bb) {
                    continue;
                }
                for (llvm::Instruction &curInst : curBB) {
                    //TODO: This might be unreliable as "dbg xxxxx" might be different!
                    //TODO: Although we set up a cache now, this can still be *slow* depending on cache hit/miss.
                    if (this->getValueStr(llvm::dyn_cast<llvm::Value>(&curInst)) == inst) {
                        return &curInst;
                    }
                }//Inst
            }//BB
        }//Func
        return nullptr;
    }

    llvm::BasicBlock *StaticAnalysisResult::getBBFromStr(std::string mod, std::string func, std::string bb) {
        //NOTE: Since now we only have one module, skip the module name match..
        for (llvm::Function &curFunc : *(this->p_module)) {
            if (curFunc.getName().str() != func) {
                continue;
            }
            for (llvm::BasicBlock &curBB : curFunc) {
                if (this->getBBStrID(&curBB) == bb) {
                    //if (curBB.getName().str() == bb) {
                    return &curBB;
                }
            }//BB
        }//Func
        return nullptr;
    }

    //TODO:
    void StaticAnalysisResult::QueryModIRsFromTagTy(std::string ty) {
        return;
    }

    std::set<uint64_t> *StaticAnalysisResult::getIoctlCmdSet(MOD_INF *p_mod_inf) {
        if (!p_mod_inf) {
            return nullptr;
        }
        std::set<uint64_t> *s = new std::set<uint64_t>();
        for (auto &x : *p_mod_inf) {
            std::set<uint64_t> &cs = x.second[1];
            s->insert(cs.begin(), cs.end());
        }
        return s;
    }

    std::string &StaticAnalysisResult::getBBStrID(llvm::BasicBlock *B) {
        static std::map<llvm::BasicBlock *, std::string> BBNameMap;
        if (BBNameMap.find(B) == BBNameMap.end()) {
            if (B) {
                if (!B->getName().empty()) {
                    BBNameMap[B] = B->getName().str();
                } else {
                    std::string Str;
                    llvm::raw_string_ostream OS(Str);
                    B->printAsOperand(OS, false);
                    BBNameMap[B] = OS.str();
                }
            } else {
                BBNameMap[B] = "";
            }
        }
        return BBNameMap[B];
    }

    //Set up a cache for the expensive "print" operation.
    std::string &StaticAnalysisResult::getValueStr(llvm::Value *v) {
        static std::map<llvm::Value *, std::string> ValueNameMap;
        if (ValueNameMap.find(v) == ValueNameMap.end()) {
            if (v) {
                std::string str;
                llvm::raw_string_ostream ss(str);
                ss << *v;
                ValueNameMap[v] = ss.str();
            } else {
                ValueNameMap[v] = "";
            }
        }
        return ValueNameMap[v];
    }

    //Set up a cache for the expensive "print" operation specifically for Type.
    std::string &StaticAnalysisResult::getTypeStr(llvm::Type *v) {
        static std::map<llvm::Type *, std::string> TypeNameMap;
        if (TypeNameMap.find(v) == TypeNameMap.end()) {
            if (v) {
                std::string str;
                llvm::raw_string_ostream ss(str);
                ss << *v;
                TypeNameMap[v] = ss.str();
            } else {
                TypeNameMap[v] = "";
            }
        }
        return TypeNameMap[v];
    }

} /* namespace sta */
