/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#include "StaticAnalysisResult.h"
#include "../DCC/general.h"

#include <iostream>

#define DEBUG_TIME 0
#define ENABLE_TAG_GROUP
#define DEBUG_TAG_GROUP 0

namespace sta {

    StaticAnalysisResult::~StaticAnalysisResult() = default;

    int StaticAnalysisResult::initStaticRes(const std::string &staticRes, dra::DataManagement *DM) {
        this->dm = DM;
        this->p_module = DM->Modules->module.get();
        try {
            std::ifstream infile;
            infile.open(staticRes);
            infile >> this->j_taintedBrs >> this->j_ctxMap >> this->j_traitMap >> this->j_tagModMap
                   >> this->j_tagConstMap >> this->j_tagInfo >> this->j_calleeMap;
            infile.close();
            this->taintedBrs = this->j_taintedBrs.get<TAINTED_BR_TY>();
            this->ctxMap = this->j_ctxMap.get<CTX_MAP_TY>();
            this->traitMap = this->j_traitMap.get<INST_TRAIT_MAP>();
            this->tagModMap = this->j_tagModMap.get<TAG_MOD_MAP_TY>();
            this->tagConstMap = this->j_tagConstMap.get<TAG_CONST_MAP_TY>();
            this->tagInfo = this->j_tagInfo.get<TAG_INFO_TY>();
            //Sort the tag info into two separate maps" global and local (e.g. user provided arg)
            for (auto &x : this->tagInfo) {
                if (x.second.find("is_global") != x.second.end() && x.second["is_global"] == "false") {
                    this->tagInfo_local[x.first] = x.second;
                } else {
                    this->tagInfo_global[x.first] = x.second;
                }
            }
            //Group the same-typed tags.
            this->setupTagGroups();
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

    llvm::DILocation *getRecursiveDILoc(llvm::Instruction *currInst, std::string &funcFileName,
                                        std::set<llvm::BasicBlock *> &visitedBBs) {
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

    //Whatever call context under which the br is tainted, we will contain its mod insts for any tags (i.e. ALL).
    MODS *StaticAnalysisResult::GetAllGlobalWriteInsts(llvm::BasicBlock *B, unsigned int branch_id) {
        BR_INF *p_taint_inf = this->QueryBranchTaint(B);
        if (!p_taint_inf) {
            std::cout << "GetAllGlobalWriteInsts : p_taint_inf = nullptr" << std::endl;
            return nullptr;
        }
        MODS *p_mod_irs = new MODS();
        //TODO: we assume now the trait for the "br" remains the same even under different contexts.
        ID_TY trait_id = 0;
        //Iterate over different contexts of "br".
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            //TODO: Consider all the tags w/ the same type.
            std::set<ID_TY> tag_ids_extend = tag_ids;
#ifdef ENABLE_TAG_GROUP
            //Consider all the tags w/ the same type.
            for (ID_TY tid : tag_ids) {
                const std::set<ID_TY> *tagGrp = this->getSameTypedTags(tid);
                if (tagGrp) {
                    tag_ids_extend.insert(tagGrp->begin(), tagGrp->end());
                }
            }
#endif
            for (ID_TY tid : tag_ids_extend) {
                //Only consider the mod insts for global taint source.
                if (this->tagInfo_local.find(tid) != this->tagInfo_local.end()) {
                    continue;
                }
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MODS *p_cur_mod_irs = this->GetRealModIrs(ps_mod_irs);

                //Append the list.
                for (auto &x : *p_cur_mod_irs) {
                    if (std::find_if(p_mod_irs->begin(), p_mod_irs->end(), [x](const Mod *m) {
                        return x->equal(m);
                    }) == p_mod_irs->end()) {
                        p_mod_irs->push_back(x);
                    }
                }
            }//tags
        }
        //According to the traits of both "br" and "store", pick out and rank the suitable mod IRs.
        //Also do some function name pair NLP analysis here.
        llvm::Instruction *inst = B->getTerminator();
        //TODO: support switch inst.
        if (llvm::dyn_cast<llvm::BranchInst>(inst)) {
            tweakModsOnTraits(p_mod_irs, trait_id, branch_id);
            filterMods(p_mod_irs, B, branch_id);
        }
        return p_mod_irs;
    }

    MODS *StaticAnalysisResult::GetAllGlobalWriteBBs(llvm::BasicBlock *B, unsigned int branch_id) {
        BR_INF *p_taint_inf = this->QueryBranchTaint(B);
        //std::cout << "p_taint_inf->size() : " << p_taint_inf->size() << std::endl;
        if (!p_taint_inf || p_taint_inf->size() == 0) {
            return nullptr;
        }
        MODS *p_mod_bbs = new MODS();
        //TODO: we assume now the trait for the "br" remains the same even under different contexts.
        ID_TY trait_id = 0;
        //Iterate over different contexts of "br".
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            std::set<ID_TY> tag_ids_extend = tag_ids;
#ifdef ENABLE_TAG_GROUP
            //Consider all the tags w/ the same type.
            for (ID_TY tid : tag_ids) {
                const std::set<ID_TY> *tagGrp = this->getSameTypedTags(tid);
                if (tagGrp) {
                    tag_ids_extend.insert(tagGrp->begin(), tagGrp->end());
                }
            }
#endif
            for (ID_TY tid : tag_ids_extend) {
                //Only consider the mod insts for global taint source.
                if (this->tagInfo_local.find(tid) != this->tagInfo_local.end()) {
                    continue;
                }
                if (this->tagModMap.find(tid) == this->tagModMap.end()) {
                    continue;
                }
                MOD_IR_TY *ps_mod_irs = &(this->tagModMap[tid]);
                MODS *p_cur_mod_bbs = this->GetRealModBbs(ps_mod_irs);

                //Append the list.
                //TODO: this can be problematic, since one BB can contain two different insts that update different global states and have different traits.
                //TODO: maybe we should deprecate GetAllGlobalWriteBBs and use GetAllGlobalWriteInsts instead.
                for (auto &x : *p_cur_mod_bbs) {
                    if (std::find_if(p_mod_bbs->begin(), p_mod_bbs->end(), [x](const Mod *m) {
                        return x->equal(m);
                    }) == p_mod_bbs->end()) {
                        p_mod_bbs->push_back(x);
                    }
                }
            }//tags
        }
        //According to the traits of both "br" and "store", pick out and rank the suitable mod IRs.
        //Also do some function name pair NLP analysis here.
        llvm::Instruction *inst = B->getTerminator();
        //TODO: support switch inst.
//        std::cout << "p_mod_bbs->size() : " << p_mod_bbs->size() << std::endl;
        if (llvm::dyn_cast<llvm::BranchInst>(inst)) {
            tweakModsOnTraits(p_mod_bbs, trait_id, branch_id);
//            std::cout << "tweakModsOnTraits p_mod_bbs->size() : " << p_mod_bbs->size() << std::endl;
            filterMods(p_mod_bbs, B, branch_id);
//            std::cout << "filterMods p_mod_bbs->size() : " << p_mod_bbs->size() << std::endl;
        }
        return p_mod_bbs;
    }

    //NOTE: this will be inclusive (the successor list also contains the root BB.)
    void StaticAnalysisResult::_get_all_successors(llvm::BasicBlock *bb, std::set<llvm::BasicBlock*> &res) {
        if (!bb || res.find(bb) != res.end()) {
            return;
        }
        //A result cache.
        if (this->succ_map.find(bb) != this->succ_map.end()) {
            res.insert(this->succ_map[bb].begin(),this->succ_map[bb].end());
            return;
        }
        //inclusive
        res.insert(bb);
        for (llvm::succ_iterator sit = llvm::succ_begin(bb), set = llvm::succ_end(bb); sit != set; ++sit) {
            this->_get_all_successors(*sit,res);
        }
        return;
    }

    //NOTE: this will be inclusive (the successor list also contains the root BB.)
    std::set<llvm::BasicBlock*> *StaticAnalysisResult::get_all_successors(llvm::BasicBlock *bb) {
        if (!bb) {
            return nullptr;
        }
        //A result cache.
        if (this->succ_map.find(bb) != this->succ_map.end()) {
            return &this->succ_map[bb];
        }
        std::set<llvm::BasicBlock*> res;
        res.clear();
        this->_get_all_successors(bb,res);
        this->succ_map[bb] = res;
        return &this->succ_map[bb];
    }

    llvm::DominatorTree *StaticAnalysisResult::get_dom_tree(llvm::Function *pfunc) {
        if (!pfunc) {
            return nullptr;
        }
        if (this->dom_map.find(pfunc) == this->dom_map.end()) {
            llvm::DominatorTree *pdom = new llvm::DominatorTree(*pfunc);
            this->dom_map[pfunc] = pdom;
        }
        return this->dom_map[pfunc];
    }

    void StaticAnalysisResult::getBranchSuccs(llvm::Instruction *inst, unsigned idx, std::set<llvm::BasicBlock*> &res) {
        if (!inst) {
            return;
        }
        llvm::BranchInst *br_inst = llvm::dyn_cast<llvm::BranchInst>(inst);
        llvm::SwitchInst *sw_inst = llvm::dyn_cast<llvm::SwitchInst>(inst);
        //TODO: This is quite awkward, from the documentation, "getNumSuccessors()" and "getSuccessor()" are member functions of "llvm::Instruction",
        //but it will throw compilation error when we try to invoke them just from "inst", maybe it's because of llvm version.
        unsigned n_succs = 0;
        if (br_inst) {
            n_succs = br_inst->getNumSuccessors();
        } else if (sw_inst) {
            n_succs = sw_inst->getNumSuccessors();
        } else {
            return;
        }
        if (idx >= n_succs) {
            return;
        }
        std::set<llvm::BasicBlock*> succ_this, succ_other;
        for (unsigned i = 0; i < n_succs; ++i) {
            llvm::BasicBlock *succ_bb = (br_inst ? br_inst->getSuccessor(i) : sw_inst->getSuccessor(i));
            std::set<llvm::BasicBlock*> *succs = this->get_all_successors(succ_bb);
            if (!succs) {
                continue;
            }
            if (i == idx) {
                succ_this.insert(succs->begin(), succs->end());
            } else {
                succ_other.insert(succs->begin(), succs->end());
            }
        }
        std::set_difference(succ_this.begin(), succ_this.end(), succ_other.begin(), succ_other.end(),
                            std::inserter(res, res.end()));
        return;
    }

    void StaticAnalysisResult::filterMods(MODS *pmods, llvm::BasicBlock *B, unsigned int branch_id) {
        if ((!pmods) || pmods->empty() || !B) {
            return;
        }
        llvm::Instruction *inst = B->getTerminator();
        if (!inst) {
            return;
        }
        //Get the successors only found for this "branch_id".
        std::set<llvm::BasicBlock*> succ_uniq;
        if (llvm::dyn_cast<llvm::BranchInst>(inst) || llvm::dyn_cast<llvm::SwitchInst>(inst)) {
            this->getBranchSuccs(inst,branch_id,succ_uniq);
        } else {
            return;
        }
        llvm::DominatorTree *pdom = this->get_dom_tree(B->getParent());
        std::remove_if(pmods->begin(), pmods->end(),
                       [succ_uniq, pdom, B](Mod *pmod) {
                           if (!pmod->B) {
                               return false;
                           }
                           //Case 0: we need to satisfy the "br" to reach the mod inst...
                           if (succ_uniq.find(pmod->B) != succ_uniq.end()) {
                               return true;
                           }
                           //Case 1: we can for sure reach the mod inst if we can reach the "br" and the mod inst is not accumulative (i.e. i++) 
                           if (pmod->B->getParent() == B->getParent() && pmod->is_trait_fixed() &&
                               pdom->dominates(pmod->B, B)) {
                               return true;
                           }
                           return false;
                       }
        );
    }

    void StaticAnalysisResult::tweakModsOnTraits(MODS *pmods, ID_TY br_trait_id, unsigned int branch_id) {
        if (!pmods) {
            return;
        }
        //TODO: verify the successor order with true/false
        bool branch = (!branch_id ? true : false);
        std::string cond("");
        int64_t v = 0;
        if (this->traitMap.find(br_trait_id) == this->traitMap.end()) {
            //No trait is also a kind of trait..
            goto CALC;
        }
        for (auto &x : this->traitMap[br_trait_id]) {
            const std::string &s = x.first;
            if (s == "==" || s == "!=") {
                if ((s == "==") == branch) {
                    //Need to take a certain value to reach the destination.
                    cond = "==";
                } else {
                    //Need to not take a certain value to reach the destination.
                    cond = "!=";
                }
                v = x.second;
            } else if (s == ">=" || s == "<=") {
                if ((s == ">=") == branch) {
                    //Need to be larger than a certain value to reach the destination.
                    cond = ">=";
                } else {
                    //Need to be smaller than a certain value to reach the destination.
                    cond = "<=";
                }
                v = x.second;
            } else if (s == ">" || s == "<") {
                if ((s == ">") == branch) {
                    //Need to be larger than a certain value to reach the destination.
                    cond = ">";
                } else {
                    //Need to be smaller than a certain value to reach the destination.
                    cond = "<";
                }
                v = x.second;
            } else if (s.substr(0, 3) == "RET") {
                //The condition is related to a function return value, do some NLP analysis.
                std::string br_func = s.substr(4);
                //E.g. if the condition is related to the return value "dequeue", then possibly to satisfy the condition we need call "enqueue" first.
                //So we need to find the "antonym" function names.
                //The heuristic is that antonym names are different but usually very similar to original names (e.g. de- and en-), so we can pick
                //those callee names with low Levenshtein distances.
                for (auto &x : this->calleeMap) {
                    int dis = this->levDistance(br_func, x.first);
                    //TODO: is "2" a proper threshold value?
                    if (dis == 0 || dis > 2) {
                        continue;
                    }
                    //Ok, we guess this is an antonym function that we should call.
                    //Get the callee instruction and treat is as a potential "Mod IR".
                    MODS *p_callee_mods = this->GetRealModBbs(&x.second);
                    if (!p_callee_mods) {
                        continue;
                    }
                    //Set proper priorities and properties of these MOD IRs.
                    for (auto &x : *p_callee_mods) {
                        x->from_nlp = true;
                    }
                    //Append these NLP Mod IRs to the original list.
                    pmods->insert(pmods->end(), p_callee_mods->begin(), p_callee_mods->end());
                }
            }
        }
        CALC:
        //Calculate mod inst priorities based given the br's and mod inst's traits.
        for (auto &x : *pmods) {
            x->calcPrio(cond, v);
        }
        //Rank the mod insts.
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
                        Mod *pmod = new Mod(pinst, &el3.second, this);
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
                        Mod *pmod = new Mod(pbb, p_mod_inf, this);
                        mod_bbs->push_back(pmod);
                    }
                }//bb
            }//func
        }//module
        return mod_bbs;
    }

    llvm::Instruction *
    StaticAnalysisResult::getInstFromStr(std::string path, std::string func, std::string bb, std::string inst) {

        auto function = this->dm->Modules->Function;
        llvm::Instruction *iii = nullptr;
        if (function.find(path) != function.end()) {
            auto file = function[path];
            if (file.find(func) != file.end()) {
                auto f = file[func];
                if (f->BasicBlock.find(bb) != f->BasicBlock.end()) {
                    auto bbb = f->BasicBlock[bb]->basicBlock;
                    for (llvm::Instruction &curInst : *bbb) {
                        if (this->getInstStrID(&curInst) == inst) {
                            iii = &curInst;
                            return iii;
                        }
                    }//Inst
                    std::cout << "not find inst : " << inst << " find bb : " << bb << std::endl;
                    bbb->dump();

                } else {
                    for (auto &it : *f->function) {
                        auto name = getBBStrID(&it);
                        if (name == bb) {
                            for (llvm::Instruction &curInst : it) {
                                if (this->getInstStrID(&curInst) == inst) {
                                    iii = &curInst;
                                    return iii;
                                }
                            }
                            std::cout << "not find inst : " << inst << std::endl;
                            it.dump();
                        }
                    }
                    std::cout << "not find bb : " << bb << std::endl;
                    f->function->dump();
                }
            } else {
                std::cout << "not find function" << std::endl;
            }
        } else {
            std::cout << "not find file : " << path << std::endl;
        }
        return iii;
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

    bool StaticAnalysisResult::isSameTypedTag(ID_TY t0, ID_TY t1) {
        if (t0 == t1) {
            return true;
        }
        if ((this->tagInfo.find(t0) == this->tagInfo.end()) != (this->tagInfo.find(t1) == this->tagInfo.end())) {
            return false;
        }
        if (this->tagInfo.find(t0) == this->tagInfo.end()) {
            //Neither exists in the map.
            return false;
        }
        auto &inf0 = this->tagInfo[t0];
        auto &inf1 = this->tagInfo[t1];
        return inf0["v"] == inf1["v"] &&
               inf0["vid"] == inf1["vid"] &&
               inf0["field"] == inf1["field"] &&
               inf0["is_global"] == inf1["is_global"] &&
               inf0["ty"] == inf1["ty"];
    }

    //Group the same typed taint tags together.
    void StaticAnalysisResult::setupTagGroups() {
#if DEBUG_TAG_GROUP
        std::cout << "-------TAG GROUP------\n";
#endif
        std::set<ID_TY> tags;
        for (auto &x : this->tagInfo) {
            tags.insert(x.first);
        }
        while (!tags.empty()) {
            ID_TY tgt = *(tags.begin());
            std::set<ID_TY> group;
            for (auto it = tags.begin(); it != tags.end(); ++it) {
                if (this->isSameTypedTag(tgt, *it)) {
                    group.insert(*it);
                    tags.erase(it);
                }
            }
#if DEBUG_TAG_GROUP
            std::cout << "+ ";
            for (auto &x : group) {
                std::cout << (const void *) x << ", ";
            }
            std::cout << "\n";
#endif
            this->tagGroups.insert(group);
        }
#if DEBUG_TAG_GROUP
        std::cout << "-------END------\n";
#endif
    }

    const std::set<ID_TY> *StaticAnalysisResult::getSameTypedTags(ID_TY tid) {
        for (auto &x : this->tagGroups) {
            if (x.find(tid) != x.end()) {
                return &x;
            }
        }
        return nullptr;
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

    std::string &StaticAnalysisResult::getInstStrID(llvm::Instruction *I) {
        static std::map<llvm::Instruction *, std::string> InstNameNoMap;
        if (InstNameNoMap.find(I) == InstNameNoMap.end()) {
            if (I) {
                if (false) {
                    //if (!I->getName().empty()){
                    InstNameNoMap[I] = I->getName().str();
                } else if (I->getParent()) {
                    int no = 0;
                    for (llvm::Instruction &i : *(I->getParent())) {
                        if (&i == I) {
                            InstNameNoMap[I] = std::to_string(no);
                            break;
                        }
                        ++no;
                    }
                } else {
                    //Seems impossible..
                    InstNameNoMap[I] = "";
                }
            } else {
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

    TRAIT *StaticAnalysisResult::getTrait(ID_TY id) {
        if (this->traitMap.find(id) != this->traitMap.end()) {
            return &(this->traitMap[id]);
        }
        return nullptr;
    }

    bool StaticAnalysisResult::getCtx(ID_TY id, std::vector<llvm::Instruction *> *pctx) {
        if (this->ctxMap.find(id) == this->ctxMap.end() || !pctx) {
            return false;
        }
        pctx->clear();
        for (auto &loc : this->ctxMap[id]) {
            llvm::Instruction *inst = this->getInstFromStr(loc[3], loc[2], loc[1], loc[0]);
            pctx->push_back(inst);
        }
        return true;
    }

    int StaticAnalysisResult::levDistance(const std::string &source, const std::string &target) {
        // Step 1
        const int n = source.length();
        const int m = target.length();
        if (n == 0) {
            return m;
        }
        if (m == 0) {
            return n;
        }

        // Good form to declare a TYPEDEF
        typedef std::vector<std::vector<int>> Tmatrix;

        Tmatrix matrix(n + 1);
        // Size the vectors in the 2.nd dimension. Unfortunately C++ doesn't
        // allow for allocation on declaration of 2.nd dimension of vec of vec
        for (int i = 0; i <= n; i++) {
            matrix[i].resize(m + 1);
        }

        // Step 2
        for (int i = 0; i <= n; i++) {
            matrix[i][0] = i;
        }
        for (int j = 0; j <= m; j++) {
            matrix[0][j] = j;
        }

        // Step 3
        for (int i = 1; i <= n; i++) {
            const char s_i = source[i - 1];

            // Step 4
            for (int j = 1; j <= m; j++) {
                const char t_j = target[j - 1];

                // Step 5
                int cost;
                if (s_i == t_j) {
                    cost = 0;
                } else {
                    cost = 1;
                }

                // Step 6
                const int above = matrix[i - 1][j];
                const int left = matrix[i][j - 1];
                const int diag = matrix[i - 1][j - 1];
                int cell = std::min(above + 1, std::min(left + 1, diag + cost));

                // Step 6A: Cover transposition, in addition to deletion,
                // insertion and substitution. This step is taken from:
                // Berghel, Hal ; Roach, David : "An Extension of Ukkonen's 
                // Enhanced Dynamic Programming ASM Algorithm"
                // (http://www.acm.org/~hlb/publications/asm/asm.html)
                if (i > 2 && j > 2) {
                    int trans = matrix[i - 2][j - 2] + 1;
                    if (source[i - 2] != t_j) trans++;
                    if (s_i != target[j - 2]) trans++;
                    if (cell > trans) cell = trans;
                }
                matrix[i][j] = cell;
            }
        }

        // Step 7
        return matrix[n][m];
    }

    //the absolute return value is the #(arg taint tags), if the value is positive, then the "br" only has arg taints,
    //if negative, there also exists global variable taints.
    int StaticAnalysisResult::getArgTaintStatus(llvm::BasicBlock *B) {
        if (!B) {
            return 0;
        }
        BR_INF *p_taint_inf = this->QueryBranchTaint(B);
        if (!p_taint_inf) {
            return 0;
        }
        bool has_global_taint = false;
        std::set<ID_TY> uniqArgTag;
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            //trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            for (ID_TY tid : tag_ids) {
                if (this->tagInfo_local.find(tid) != this->tagInfo_local.end()) {
                    uniqArgTag.insert(tid);
                } else if (this->tagInfo_global.find(tid) != this->tagInfo_global.end()) {
                    has_global_taint = true;
                }
            }
        }
        int n = uniqArgTag.size();
        if (has_global_taint) {
            n = 0 - n;
        }
        return n;
    }

    bool StaticAnalysisResult::getAllTagConstants(ID_TY tag_id, CONST_INF *p_consts) {
        if (this->tagConstMap.find(tag_id) == this->tagConstMap.end() || !p_consts) {
            return false;
        }
        for (auto &i0 : this->tagConstMap[tag_id]) {
            //i0.first : file
            for (auto &i1 : i0.second) {
                //i1.first : func
                for (auto &i2 : i1.second) {
                    //i2.first : BB
                    for (auto &i3 : i2.second) {
                        //i3.first : inst
                        p_consts->insert(i3.second.begin(), i3.second.end());
                    }
                }
            }
        }
        return true;
    }

    //parse the type str to get more easy-to-read info.
    //type str, e.g.
    //%struct.A = type {...}:3->%struct.B = type {...}:2.%struct.C = type {...}:5
    std::vector<FieldPtr *> *StaticAnalysisResult::parseTypeStr(std::string tys) {
        std::vector<FieldPtr *> *r = new std::vector<FieldPtr *>();
        int prev = 0;
        for (int i = 0; i < tys.length(); ++i) {
            if (tys[i] == ':') {
                std::string obj_ty = tys.substr(prev, i - prev);
                size_t j = 0;
                long field = 0;
                try {
                    field = std::stoi(tys.substr(i + 1), &j, 10);
                } catch (...) {
                    std::cout << "Exception in StaticAnalysisResult::parseTypeStr: " << tys << "\n";
                    return r;
                }
                if (j <= 0) {
                    std::cout << "No valid field number after :" << tys << "\n";
                    return r;
                }
                i += (j + 1);
                bool is_embed = true;
                if (i < tys.length() - 1) {
                    if (tys[i] == '-' && tys[i + 1] == '>') {
                        is_embed = false;
                        prev = i + 2;
                    } else if (tys[i] != '.') {
                        //Something goes wrong
                        std::cout << "Invalid format :" << tys << "\n";
                        return r;
                    } else {
                        prev = i + 1;
                    }
                }
                FieldPtr *fiptr = new FieldPtr(obj_ty, field, is_embed);
                r->push_back(fiptr);
            }
        }
        return r;
    }

    //Return the type of the tag variable, see comment
    std::vector<std::vector<FieldPtr *> *> *StaticAnalysisResult::getTagType(ID_TY tag_id) {
        if (this->tagInfo.find(tag_id) == this->tagInfo.end()) {
            return nullptr;
        }
        std::vector<std::vector<FieldPtr *> *> *r = new std::vector<std::vector<FieldPtr *> *>();
        for (auto &x : this->tagInfo[tag_id]) {
            if (x.first.find("hs_") == 0) {
                //Find a hierarchy string.
                r->push_back(this->parseTypeStr(x.second));
            }
        }
        return r;
    }

    //If the conditional jump in the passed-in BB is tainted by any user arg,
    //return these args (e.g. a field in a user struct) and related constants we collected.
    //NOTE: free the returned obj after using.
    std::map<ID_TY, CONST_INF> *StaticAnalysisResult::getArgTaintInfo(llvm::BasicBlock *B) {
        if (!B) {
            return nullptr;
        }
        BR_INF *p_taint_inf = this->QueryBranchTaint(B);
        if (!p_taint_inf) {
            return nullptr;
        }
        std::map<ID_TY, CONST_INF> *pres = new std::map<ID_TY, CONST_INF>();
        for (auto &x : *p_taint_inf) {
            auto &actx_id = x.first;
            //trait_id = std::get<0>(x.second);
            auto &tag_ids = std::get<1>(x.second);
            for (ID_TY tid : tag_ids) {
                if (this->tagInfo_local.find(tid) != this->tagInfo_local.end()) {
                    //Ok, find one arg that taints this BB.
                    //Get its collected cmp constants.
                    this->getAllTagConstants(tid, &((*pres)[tid]));
                }
            }
        }
        return pres;
    }

} /* namespace sta */
