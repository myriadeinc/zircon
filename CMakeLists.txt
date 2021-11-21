cmake_minimum_required (VERSION 2.8.11)
project(zircon)

if (DEFINED ENV{MONERO_DIR})
    get_filename_component(MONERO_DIR $ENV{MONERO_DIR} ABSOLUTE)
    message("Using Monero source from env ${MONERO_DIR}")
else()
    get_filename_component(MONERO_DIR "${CMAKE_SOURCE_DIR}/../xmrsrc" ABSOLUTE)
    message("Monero surce directory is not defined, using default ${MONERO_DIR}")
endif()

add_subdirectory(xmrlib)

add_custom_command(
    OUTPUT build/bin
    COMMAND build/env.sh go get -v ./...
)

add_custom_target(build ALL DEPENDS xmrlib ${CMAKE_CURRENT_BINARY_DIR}/build/bin)
