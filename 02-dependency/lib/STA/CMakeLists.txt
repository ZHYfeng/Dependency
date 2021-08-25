set(SOURCES
        StaticAnalysisResult.cpp
        )

add_library(STA ${SOURCES})
set_target_properties(STA PROPERTIES LINKER_LANGUAGE CXX)
target_link_libraries(STA PUBLIC
        DCC
        )