/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#include "StaticAnalysisResult.h"

#include <iostream>

#define DEBUG_TIME 0

namespace sta {

    StaticAnalysisResult::~StaticAnalysisResult() = default;

    int StaticAnalysisResult::initStaticRes(const std::string &staticRes, dra::DataManagement *DM) {
        this->dm = DM;
        this->p_module = DM->Modules->module.get();
        try {
            std::ifstream infile;
            infile.open(staticRes);
            infile >> this->j_taintedBrs >> this->j_ctxMap >> this->j_traitMap >> this->j_tagModMap >> this->j_tagInfo >> this->j_calleeMap;
            infile.close();
            this->taintedBrs = this->j_taintedBrs.get<TAINTED_BR_TY>();
            this->ctxMap = this->j_ctxMap.get<CTX_MAP_TY>();
            this->traitMap = this->j_traitMap.get<INST_TRAIT_MAP>();
            this->tagModMap = this->j_tagModMap.get<TAG_MOD_MAP_TY>();
            this->tagInfo = this->j_tagInfo.get<TAG_INFO_TY>();
            this->calleeMap = this->j_calleeMap.get<CALLEE_MAP_TY>();
            return 0;
        } catch (...) {
            std::cout << "Fail to deserialize the static analysis results!\n";
        }
        return 1;
    }

    void StaticAnalysisResult::stripFuncNameSuffix(std::string *fn) {
        if (!fn) {
            return;
        }
        std::size_t n = fn->rfind(".");
        if (n != std::string::npos) {
            fn->erase(n);
        }
        return;
    }

    std::string getFunctionFileName(llvm::Function *F) {
        llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
        F->getAllMetadata(MDs);
        for (auto &MD : MDs) {
            if (llvm::MDNode *N = MD.second) {
                if (auto *subProgram = llvm::dyn_cast<llvm::DISubprogram>(N)) {
                    return subProgram->getFilename();
                }
            }
        }
        return "";
    }

    llvm::DILocation *getRecursiveDILoc(llvm::Instruction *currInst, std::string &funcFileName, std::set<llvm::BasicBlock *> &visitedBBs) {
        llvm::DILocation *currIL = currInst->getDebugLoc().get();
        if (funcFileName.length() == 0) {
            return currIL;
        }
        if (currIL != nullptr && currIL->getFilename().equals(funcFileName)) {
            return currIL;
        }

        llvm::BasicBlock *currBB = currInst->getParent();
        if (visitedBBs.find(currBB) != visitedBBs.end()) {
            return nullptr;
        }
        for (auto &iu :currBB->getInstList()) {
            llvm::Instruction *currIterI = &iu;
            llvm::DILocation *currIteDL = currIterI->getDebugLoc();
            if (currIteDL != nullptr && currIteDL->getFilename().equals(funcFileName)) {
                return currIteDL;
            }
            if (currIterI == currInst) {
                break;
            }
        }

        visitedBBs.insert(currBB);


        for (auto it = llvm::pred_begin(currBB), et = llvm::pred_end(currBB); it != et; ++it) {
            llvm::BasicBlock *predecessor = *it;
            llvm::DILocation *currBBLoc = getRecursiveDILoc(predecessor->getTerminator(), funcFileName, visitedBBs);
            if (currBBLoc != nullptr) {
                return currBBLoc;
            }
        }
        return nullptr;
    }

    llvm::DILocation *StaticAnalysisResult::getCorrectInstrLocation(llvm::Instruction *I) {
        llvm::DILocation *instrLoc = I->getDebugLoc().get();
        //BasicBlock *firstBB = &(I->getFunction()->getEntryBlock());
        //Instruction *firstInstr = firstBB->getFirstNonPHIOrDbg();

        //DILocation *firstIL = firstInstr->getDebugLoc().get();
        std::set<llvm::BasicBlock *> visitedBBs;
        std::string funcFileName = getFunctionFileName(I->getFunction());


        if (instrLoc != nullptr && instrLoc->getFilename().endswith(".c")) {
            return instrLoc;
        }

        if (instrLoc == nullptr || (funcFileName.length() > 0 && !instrLoc->getFilename().equals(funcFileName))) {
            // OK, the instruction is from the inlined function.
            visitedBBs.clear();
            llvm::DILocation *actualLoc = getRecursiveDILoc(I, funcFileName, visitedBBs);
            if (actualLoc != nullptr) {

                return actualLoc;
            }
        }

        return instrLoc;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(llvm::Instruction *I, bool skip_inst) {
        if (!I) {
            return nullptr;
        }
        std::string inst(""), bb, func, file;
        if (!skip_inst) {
            std::cout << "getLocInf : !skip_inst" << std::endl;
            inst = this->getInstStrID(I);
        }
        llvm::DILocation *instrLoc = StaticAnalysisResult::getCorrectInstrLocation(I);
        if (I->getParent()) {
            bb = this->getBBStrID(I->getParent());
        }
        if (I->getFunction()) {
            func = I->getFunction()->getName().str();
            this->stripFuncNameSuffix(&func);
        }
        //Put the file name.
        if (instrLoc != nullptr) {
            file = instrLoc->getFilename();
        } else {
            //TODO: not sure what to do here..
            if (I->getModule()) {
                file = I->getModule()->getName().str();
            } else {
                //Is this possible?
            }
        }
        LOC_INF *str_inst = new LOC_INF;
        str_inst->push_back(inst);
        str_inst->push_back(bb);
        str_inst->push_back(func);
        str_inst->push_back(file);
        return str_inst;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(llvm::BasicBlock *B) {
        if (!B) {
            std::cout << "getLocInf : b = nullptr" << std::endl;
            return nullptr;
        }
        return this->getLocInf(&*(B->begin()), true);
    }

    //Given a bb, return the taint information regarding its last br inst.
    //The returned info is a map from the context id to the taint tag id set.
    BR_INF *StaticAnalysisResult::QueryBranchTaint(llvm::BasicBlock *B) {
        if (!B) {
            std::cout << "QueryBranchTaint : b = bullptr" << std::endl;
            return nullptr;
        }
        LOC_INF *p_loc = this->getLocInf(B);
        if (!p_loc) {
            std::cout << "QueryBranchTaint : p_loc = nullptr" << std::endl;
            return nullptr;
        }
        auto &res3 = this->taintedBrs;
        if (res3.find((*p_loc)[3]) != res3.end()) {
//            std::cout << "QueryBranchTaint : find path : " << (*p_loc)[3] << std::endl;
            auto &res2 = res3[(*p_loc)[3]];
            if (res2.find((*p_loc)[2]) != res2.end()) {
//                std::cout << "QueryBranchTaint : find function : " << (*p_loc)[2] << std::endl;
                auto &res1 = res2[(*p_loc)[2]];
                if (res1.find((*p_loc)[1]) != res1.end()) {
//                    std::cout << "QueryBranchTaint : find bb : " << (*p_loc)[1] << std::endl;
                    return &(res1[(*p_loc)[1]]);
                }
            }
        }
//        std::cout << "QueryBranchTaint : return = nullptr : out side analysis" << std::endl;
        return nullptr;
    }

    MODS *StaticAnalysisResult::GetAllGlobalWriteInsts(llvm::BasicBlock *B) {
        return this->GetAllGlobalWriteInsts(this->QueryBranchTaint(B));
    }

    //Whatever call context under which the br is tainted, we will contain its mod insts for any tags (i.e. ALL).
    MODS *StaticAnalysisResult::GetAllGlobalWriteInsts(BR_INF *p_taint_inf) {
        if (!p_taint_inf) {
            std::cout << "GetAllGlobalWriteInsts : p_taint_inf = nullptr" << std::endl;
            return nullptr;
        }
        MODS *p_mod_irs = new MODS();
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            auto &trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            for (ID_TY tid : tag_ids) {
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MODS *p_cur_mod_irs = this->GetRealModIrs(ps_mod_irs);

                //Append the list.
                p_mod_irs->insert(p_mod_irs->end(),p_cur_mod_irs->begin(),p_cur_mod_irs->end());
            }//tags
        }
        return p_mod_irs;
    }

    MODS *StaticAnalysisResult::GetAllGlobalWriteBBs(llvm::BasicBlock *B) {
        return this->GetAllGlobalWriteBBs(this->QueryBranchTaint(B));
    }

    MODS *StaticAnalysisResult::GetAllGlobalWriteBBs(BR_INF *p_taint_inf) {
        if (!p_taint_inf) {
            return nullptr;
        }
        MODS *p_mod_bbs = new MODS();
        //TODO: we assume now the trait for the "br" remains the same even under different contexts.
        ID_TY trait_id = 0;
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            for (ID_TY tid : tag_ids) {
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MODS *p_cur_mod_bbs = this->GetRealModBbs(ps_mod_irs);

                //Append the list.
                p_mod_bbs->insert(p_mod_bbs->end(),p_cur_mod_bbs->begin(),p_cur_mod_bbs->end());
            }//tags
        }
        //According to the traits of both "br" and "store", pick out and rank the suitable mod IRs.
        //Also do some function name pair NLP analysis here.
        //TODO: the "bool" value indicating the taken branch.
        tweakModsOnTraits(p_mod_bbs,trait_id,true);
        return p_mod_bbs;
    }

    void StaticAnalysisResult::tweakModsOnTraits(MODS *pmods, ID_TY br_trait_id, bool branch) {
        if ((!pmods) || this->traitMap.find(br_trait_id) == this->traitMap.end()) {
            return;
        }
        TRAIT& br_trait = this->traitMap[br_trait_id];
        for (auto& x : br_trait) {
            if (x.first == "==" || x.first == "!=") {
                if ((x.first == "==") == branch) {
                    //Need to take a certain value to reach the destination.
                }else {
                    //Need to not take a certain value to reach the destination.
                }
            }else if (x.first == ">=" || x.first == "<=") {
                if ((x.first == ">=") == branch) {
                    //Need to be larger than a certain value to reach the destination.
                }else {
                    //Need to be smaller than a certain value to reach the destination.
                }
            }else if (x.first.substr(0,3) == "RET") {
                //The condition is related to a function return value, do some NLP analysis.
            }
        }
    }

    MODS *StaticAnalysisResult::GetRealModIrs(MOD_IR_TY *p_mod_irs) {
        if (!p_mod_irs) {
            return nullptr;
        }
        MODS *mod_irs = new MODS();
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
                        Mod *pmod = new Mod(pinst,&el3.second);
                        mod_irs->push_back(pmod);
                    }//inst
                }//bb
            }//func
        }//module
        return mod_irs;
    }

    MODS *StaticAnalysisResult::GetRealModBbs(MOD_IR_TY *p_mod_irs) {
        if (!p_mod_irs) {
            return nullptr;
        }
        MODS *mod_bbs = new MODS();
        for (auto &el0 : *p_mod_irs) {
            const std::string &path = el0.first;
            for (auto &el1 : (*p_mod_irs)[path]) {
                const std::string &func = el1.first;
                for (auto &el2 : (*p_mod_irs)[path][func]) {
                    const std::string &bb = el2.first;
                    llvm::BasicBlock *pbb = this->getBBFromStr(path, func, bb);
                    if (!pbb) {
                        continue;
                    }
                    //We should choose the last inst in the bb since it will overwrite the previous written value.
                    //We can easily do this since we now use the inst's loop# within the parent BB to identify it.
                    int k = -1;
                    MOD_INF *p_mod_inf = nullptr;
                    for (auto &el3 : (*p_mod_irs)[path][func][bb]) {
                        if (k < 0 || std::stoi(el3.first) > k) {
                            k = std::stoi(el3.first);
                            p_mod_inf = &el3.second;
                        }
                    }
                    if (p_mod_inf) {
                        Mod *pmod = new Mod(pbb,p_mod_inf);
                        mod_bbs->push_back(pmod);
                    }
                }//bb
            }//func
        }//module
        return mod_bbs;
    }

    llvm::Instruction *StaticAnalysisResult::getInstFromStr(std::string path, std::string func, std::string bb, std::string inst) {

        auto function = this->dm->Modules->Function;
        if (function.find(path) != function.end()) {
            auto file = function[path];
            if (file.find(func) != file.end()) {
                auto f = file[func];
                if (f->BasicBlock.find(bb) != f->BasicBlock.end()) {
                    auto bbb = f->BasicBlock[bb]->basicBlock;
                    for (llvm::Instruction &curInst : *bbb) {
                        if (this->getInstStrID(&curInst) == inst) {
                            return &curInst;
                        }
                    }//Inst
                } else {
                    for (auto &it : *f->function) {
                        auto name = getBBStrID(&it);
                        if (name == bb) {
                            for (llvm::Instruction &curInst : it) {
                                if (this->getInstStrID(&curInst) == inst) {
                                    return &curInst;
                                }
                            }
                        }
                    }
                }
            }
        }
        return nullptr;
    }

    llvm::BasicBlock *StaticAnalysisResult::getBBFromStr(std::string path, std::string func, std::string bb) {
#if DEBUG_TIME
        std::time_t current_time = std::time(NULL);
        std::cout << std::ctime(&current_time) << "*time : getBBFromStr" << std::endl;
#endif

        llvm::BasicBlock *bbb = nullptr;
        auto function = this->dm->Modules->Function;
        if (function.find(path) != function.end()) {
            auto file = function[path];
            if (file.find(func) != file.end()) {
#if DEBUG_TIME
                current_time = std::time(NULL);
                std::cout << std::ctime(&current_time) << "*time : getBBFromStr function" << std::endl;
#endif
                auto f = file[func];
                if (f->BasicBlock.find(bb) != f->BasicBlock.end()) {
                    bbb = f->BasicBlock[bb]->basicBlock;
#if DEBUG_TIME
                    current_time = std::time(NULL);
                    std::cout << std::ctime(&current_time) << "*time : getBBFromStr basicBlock" << std::endl;
#endif
                    return bbb;
                } else {
#if DEBUG_TIME

                    current_time = std::time(NULL);
                    std::cout << std::ctime(&current_time) << "*time : getBBFromStr for (auto &it : *f->function) {" << std::endl;
#endif
                    for (auto &it : *f->function) {
                        auto name = getBBStrID(&it);
                        if (name == bb) {
                            bbb = &it;
                            return bbb;
                        }
                    }
#if DEBUG_TIME

                    current_time = std::time(NULL);
                    std::cout << std::ctime(&current_time) << "*time : getBBFromStr finish for (auto &it : *f->function) {" << std::endl;
#endif
                }
            } else {
                std::cout << "not find function" << std::endl;
            }
        } else {
            std::cout << "not find file" << std::endl;
        }
        return bbb;
    }

    //TODO:
    void StaticAnalysisResult::QueryModIRsFromTagTy(std::string ty) {
        return;
    }

    /*
    std::string &StaticAnalysisResult::getBBStrID(llvm::BasicBlock *B) {
        std::time_t current_time = std::time(NULL);
        std::cout << std::ctime(&current_time) << "*time : start getBBStrID" << std::endl;

        static std::map<llvm::BasicBlock *, std::string> BBNameMap;
        if (BBNameMap.find(B) == BBNameMap.end()) {
            if (B) {
                if (!B->getName().empty()) {
                    BBNameMap[B] = B->getName().str();
                } else {
                    std::string Str;
                    llvm::raw_string_ostream OS(Str);
                    current_time = std::time(NULL);
                    std::cout << std::ctime(&current_time) << "*time : start printAsOperand" << std::endl;
                    B->printAsOperand(OS, false);
                    current_time = std::time(NULL);
                    std::cout << std::ctime(&current_time) << "*time : finish printAsOperand" << std::endl;
                    BBNameMap[B] = OS.str();
                }
            } else {
                BBNameMap[B] = "";
            }
        }
        current_time = std::time(NULL);
        std::cout << std::ctime(&current_time) << "*time : finish getBBStrID" << BBNameMap[B] << std::endl;
        return BBNameMap[B];
    }
    */

    std::string &StaticAnalysisResult::getBBStrID(llvm::BasicBlock *B) {
        static std::map<llvm::BasicBlock *, std::string> BBNameMap;
        if (BBNameMap.find(B) == BBNameMap.end()) {
            if (B) {
                if (!B->getName().empty()) {
                    BBNameMap[B] = B->getName().str();
                } else if (B->getParent()) {
                    int no = 0;
                    for (llvm::BasicBlock &bb : *(B->getParent())) {
                        if (&bb == B) {
                            BBNameMap[B] = std::to_string(no);
                            break;
                        }
                        ++no;
                    }
                } else {
                    //Seems impossible..
                    BBNameMap[B] = "";
                }
            } else {
                BBNameMap[B] = "";
            }
        }
        return BBNameMap[B];
    }

    std::string& StaticAnalysisResult::getInstStrID(llvm::Instruction* I) {
        static std::map<llvm::Instruction*,std::string> InstNameNoMap;
        if (InstNameNoMap.find(I) == InstNameNoMap.end()) {
            if (I) {
                if (false){
                //if (!I->getName().empty()){
                    InstNameNoMap[I] = I->getName().str();
                }else if (I->getParent()){
                    int no = 0;
                    for (llvm::Instruction& i : *(I->getParent())) {
                        if (&i == I) {
                            InstNameNoMap[I] = std::to_string(no);
                            break;
                        }
                        ++no;
                    }
                }else{
                    //Seems impossible..
                    InstNameNoMap[I] = "";
                }
            }else{
                InstNameNoMap[I] = "";
            }
        }
        return InstNameNoMap[I];
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
