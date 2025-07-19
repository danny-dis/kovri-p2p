
add_library(
    kovri-core-riffle
    protocol.cc)

target_include_directories(kovri-core-riffle PUBLIC ${CMAKE_SOURCE_DIR}/src ${CMAKE_SOURCE_DIR}/deps)

target_link_libraries(
    kovri-core-riffle
    kovri-core-crypto
    cryptopp)
