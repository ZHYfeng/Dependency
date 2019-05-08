/*
* Class to deserialize and query the static analysis results.
* By hz
* 05/08/2019
*/

#include "StaticAnalysisResult.h"

namespace STA {

    StaticAnalysisResult::~StaticAnalysisResult() = default;

    int StaticAnalysisResult::initStaticRes(const std::string &staticRes, llvm::Module *p_module) {
        this->p_module = p_module;
        try{
            std::ifstream infile;
            infile.open(staticRes);
            infile >> this->j_taintedBrs >> this->j_analysisCtxMap >> this->j_tagMap >> this->j_modInstCtxMap;
            infile.close();
            return 0;
        } catch (...) {
            std::cout << "Fail to deserialize the static analysis results!\n";
        }
        return 1;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(llvm::Instruction* I) {
        if(!I){
            return nullptr;
        }
        std::string inst,bb,func,mod;
        std::string str;
        llvm::raw_string_ostream ss(str);
        ss << *I;
        inst = ss.str();
        if(I->getParent()){
            bb = I->getParent()->getName().str();
        }
        if(I->getFunction()){
            func = I->getFunction()->getName().str();
        }
        if(I->getModule()){
            mod = I->getModule()->getName().str();
        }
        LOC_INF *loc_inf = new LOC_INF;
        loc_inf->push_back(inst);
        loc_inf->push_back(bb);
        loc_inf->push_back(func);
        loc_inf->push_back(mod);
        return loc_inf;
    }

    LOC_INF *StaticAnalysisResult::getLocInf(llvm::BasicBlock* B) {
        if(!B){
            return nullptr;
        }
        return this->getLocInf(&*(B->begin()));
    }

    nlohmann::json *StaticAnalysisResult::findLocInJson(LOC_INF *p_loc, unsigned int e, nlohmann::json *data) {
        if ((!p_loc) || p_loc->size() < e || !data){
            return nullptr;
        }
        nlohmann::json *res = data;
        for (unsigned int i = 3; i >= e; --i){
            std::string k = (*p_loc)[i];
            if (res->find(k) == res->end()){
                return nullptr;
            }
            res = &((*res)[k]);
        }
        return res;
    }

    //Given a bb, return the taint information regarding its last br inst.
    //The returned info is a map from the context id to the taint tag id set.
    nlohmann::json *StaticAnalysisResult::QueryBranchTaint(llvm::BasicBlock* B) {
        if (!B) {
            return nullptr;
        }
        LOC_INF *p_loc = this->getLocInf(B);
        if (!p_loc){
            return nullptr;
        }
        nlohmann::json *pj_taint_inf = this->findLocInJson(p_loc,1,this->j_taintedBrs);
        if(!pj_taint_inf){
            //This means the br instruction of this bb is not tainted by global states. 
            return nullptr;
        }
        return pj_taint_inf;
    }

    MOD_IRS *StaticAnalysisResult::GetAllGlobalWriteInsts(llvm::BasicBlock* B) {
        return this->GetAllGlobalWriteInsts(this->QueryBranchTaint(B));
    }

    //Whatever call context under which the br is tainted, we will contain its mod insts for any tags (i.e. ALL).
    MOD_IRS *StaticAnalysisResult::GetAllGlobalWriteInsts(nlohmann::json *pj_taint_inf) {
        if (!pj_taint_inf){
            return nullptr;
        }
        MOD_IRS *p_mod_irs new MOD_IRS();
        for (auto& el : pj_taint_inf->items()) {
            //Analysis context id under which this br is tainted.
            nlohmann::json j_actx_id = el.key();
            //A set of taint tag ids for this tainted br.
            nlohmann::json j_tag_ids = el.value();
            for (auto& tid : j_tag_ids) {
                nlohmann::json *pj_mod_irs = this->QueryModIRsFromTagID(tid.get<unsigned long>());
                if (!pj_mod_irs){
                    continue;
                }
                //Get the mod insts for current taint tag, note that we may have multiple tags for one tainted br...
                MOD_IRS *p_cur_mod_irs = this->j2ModIrs(pj_mod_irs);
                //Merge.
                for (auto const& x : *p_cur_mod_irs) {
                    if (p_mod_irs->find(x.first) != p_mod_irs->end()) {
                        (*p_mod_irs)[x.first].insert(x.second.begin(),x.second.end());
                    }else{
                        (*p_mod_irs)[x.first] = x.second;
                    }
                }//merge
            }//tags
        }
        return p_mod_irs;
    }

    MOD_BBS *StaticAnalysisResult::GetAllGlobalWriteBBs(llvm::BasicBlock* B) {
        return this->GetAllGlobalWriteBBs(this->QueryBranchTaint(B));
    }

    MOD_BBS *StaticAnalysisResult::GetAllGlobalWriteBBs(nlohmann::json *pj_taint_inf) {
        if (!pj_taint_inf){
            return nullptr;
        }
        MOD_BBS *p_mod_bbs new MOD_BBS();
        for (auto& el : pj_taint_inf->items()) {
            //Analysis context id under which this br is tainted.
            nlohmann::json j_actx_id = el.key();
            //A set of taint tag ids for this tainted br.
            nlohmann::json j_tag_ids = el.value();
            for (auto& tid : j_tag_ids) {
                nlohmann::json *pj_mod_irs = this->QueryModIRsFromTagID(tid.get<unsigned long>());
                if (!pj_mod_irs){
                    continue;
                }
                //Get the mod insts for current taint tag, note that we may have multiple tags for one tainted br...
                MOD_BBS *p_cur_mod_bbs = this->j2ModBbs(pj_mod_irs);
                //Merge.
                for (auto const& x : *p_cur_mod_bbs) {
                    if (p_mod_bbs->find(x.first) != p_mod_bbs->end()) {
                        (*p_mod_bbs)[x.first].insert(x.second.begin(),x.second.end());
                    }else{
                        (*p_mod_bbs)[x.first] = x.second;
                    }
                }//merge
            }//tags
        }
        return p_mod_bbs;
    }

    MOD_IRS *StaticAnalysisResult::j2ModIrs(nlohmann::json *pj_mod_irs) {
        if(!pj_mod_irs) {
            return nullptr;
        }
        MOD_IRS *mod_irs = new MOD_IRS();
        for (auto& el0 : pj_mod_irs->items()) {
            std::string module = el0.key().get<std::string>();
            for (auto& el1 : (*pj_mod_irs)[module].items()) {
                std::string func = el1.key().get<std::string>();
                for (auto& el2 : (*pj_mod_irs)[module][func].items()) {
                    std::string bb = el2.key().get<std::string>();
                    for (auto& el3 : (*pj_mod_irs)[module][func][bb].items()) {
                        std::string inst = el3.key().get<std::string>();
                        //Get the actual Instruction* according to these string info
                        llvm::Instruction *pinst = this->getInstFromStr(module,func,bb,inst);
                        if (!pinst) {
                            continue;
                        }
                        (*mod_irs)[pinst] = el3.value().get<MOD_INF>();
                    }//inst
                }//bb
            }//func
        }//module
        return mod_irs;
    }

    MOD_BBS *StaticAnalysisResult::j2ModBbs(nlohmann::json *pj_mod_irs) {
        if(!pj_mod_irs) {
            return nullptr;
        }
        MOD_BBS *mod_bbs = new MOD_BBS();
        for (auto& el0 : pj_mod_irs->items()) {
            std::string module = el0.key().get<std::string>();
            for (auto& el1 : (*pj_mod_irs)[module].items()) {
                std::string func = el1.key().get<std::string>();
                for (auto& el2 : (*pj_mod_irs)[module][func].items()) {
                    std::string bb = el2.key().get<std::string>();
                    llvm::BasicBlock *pbb = this->getBBFromStr(module,func,bb);
                    if (!pbb) {
                        continue;
                    }
                    for (auto& el3 : (*pj_mod_irs)[module][func][bb].items()) {
                        MOD_INF& mod_inf = el3.value().get<MOD_INF>();
                        (*mod_bbs)[pbb].insert(mod_inf.begin(),mod_inf.end());
                    }//inst
                }//bb
            }//func
        }//module
        return mod_bbs;
    }

    llvm::Instruction *StaticAnalysisResult::getInstFromStr(std::string mod, std::string func, std::string bb, std::string inst) {
        //NOTE: Since now we only have one module, skip the module name match..
        for (llvm::Function& curFunc : *(this->p_module)) {
            if (curFunc.getName().str() != func){
                continue;
            }
            for (llvm::BasicBlock& curBB : curFunc) {
                if (curBB.getName().str() != bb) {
                    continue;
                }
                for (llvm::Instruction& curInst : curBB) {
                    //TODO: This might be unreliable as "dbg xxxxx" might be different!
                    //TODO: This can be *slow* since dump llvm::Instruction is time-consuming!
                    llvm::raw_string_ostream ss(str);
                    ss << curInst;
                    if (ss.str() == inst) {
                        return &curInst;
                    }
                }//Inst
            }//BB
        }//Func
        return nullptr;
    }

    llvm::BasicBlock *StaticAnalysisResult::getBBFromStr(std::string mod, std::string func, std::string bb) {
        //NOTE: Since now we only have one module, skip the module name match..
        for (llvm::Function& curFunc : *(this->p_module)) {
            if (curFunc.getName().str() != func){
                continue;
            }
            for (llvm::BasicBlock& curBB : curFunc) {
                if (curBB.getName().str() == bb) {
                    return &curBB;
                }
            }//BB
        }//Func
        return nullptr;
    }

    nlohmann::json *StaticAnalysisResult::QueryModIRsFromTagID(unsigned long tid) {
        if (!this->j_tagMap){
            return nullptr;
        }
        if (this->j_tagMap.find(tid) == this->j_tagMap.end()) {
            //No taint tag with the specified id.
            return nullptr;
        }
        for (auto& el : this->j_tagMap[tid].items()) {
            return &(el.value());
        }
        return nullptr;
    }

    //TODO:
    nlohmann::json *StaticAnalysisResult::QueryModIRsFromTagTy(std::string ty) {
        if (!this->j_tagMap){
            return nullptr;
        }
        for (auto& el : this->j_tagMap.items()) {
            nlohmann::json *pj = &(el.value());
            if (pj->find(ty) != pj->end()){
                //TODO: merge from multiple tags w/ the same type.
                //return &((*pj)[ty]);
            }
        }
        return nullptr;
    }

    std::set<uint64_t> *StaticAnalysisResult::getIoctlCmdSet(MOD_INF* p_mod_inf) {
        if (!p_mod_inf) {
            return nullptr;
        }
        std::set<uint64_t> *s = new std::set<uint64_t>();
        for (auto& x : *p_mod_inf) {
            std::set<uint64_t>& cs = x.second[1];
            s->insert(cs.begin(),cs.end());
        }
        return s;
    }

} /* namespace STA */
