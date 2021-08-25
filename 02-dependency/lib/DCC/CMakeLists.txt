set(SOURCES
        DependencyControlCenter.cpp
        general.cpp)

add_library(DCC ${SOURCES})
set_target_properties(DCC PROPERTIES LINKER_LANGUAGE CXX)
target_link_libraries(DCC PUBLIC
        DMM
        )