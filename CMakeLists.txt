cmake_minimum_required(VERSION 3.20)

project(NextCmd LANGUAGES NONE)

find_program(GO_EXECUTABLE NAMES go REQUIRED)

option(
    NEXTCMD_OUTPUT_IN_SOURCE_ROOT
    "Write the default NextCmd executable to the repository root"
    OFF
)

set(NEXTCMD_ARTIFACT_DIR "${CMAKE_BINARY_DIR}/bin")
set(
    NEXTCMD_GO_CACHE_DIR
    "${CMAKE_BINARY_DIR}/go-cache"
    CACHE PATH
    "Go build cache used by CMake targets"
)
if(WIN32)
    set(NEXTCMD_EXECUTABLE_NAME "NextCmd.exe")
else()
    set(NEXTCMD_EXECUTABLE_NAME "NextCmd")
endif()

# Keep normal builds out of the source tree. The opt-in mode is useful for
# manually testing NextCmd from the repository root.
if(NEXTCMD_OUTPUT_IN_SOURCE_ROOT)
    set(NEXTCMD_NATIVE_OUTPUT_DIR "${CMAKE_SOURCE_DIR}")
else()
    set(NEXTCMD_NATIVE_OUTPUT_DIR "${NEXTCMD_ARTIFACT_DIR}")
endif()
set(NEXTCMD_NATIVE_OUTPUT "${NEXTCMD_NATIVE_OUTPUT_DIR}/${NEXTCMD_EXECUTABLE_NAME}")

# CONFIGURE_DEPENDS makes CMake reconfigure when a Go source file is added or
# removed. Listing only project source locations avoids scanning build trees.
file(GLOB_RECURSE NEXTCMD_GO_SOURCES CONFIGURE_DEPENDS
    "${CMAKE_SOURCE_DIR}/cmd/*.go"
    "${CMAKE_SOURCE_DIR}/examples/*.go"
    "${CMAKE_SOURCE_DIR}/internal/*.go"
    "${CMAKE_SOURCE_DIR}/plugins/*.go"
    "${CMAKE_SOURCE_DIR}/sdk/*.go"
)
list(APPEND NEXTCMD_GO_SOURCES
    "${CMAKE_SOURCE_DIR}/main.go"
    "${CMAKE_SOURCE_DIR}/go.mod"
)

# Model the executable as a generated output so unchanged builds are skipped.
add_custom_command(
    OUTPUT "${NEXTCMD_NATIVE_OUTPUT}"
    COMMAND "${CMAKE_COMMAND}" -E make_directory "${NEXTCMD_NATIVE_OUTPUT_DIR}"
    COMMAND "${CMAKE_COMMAND}" -E env
        "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
        "${GO_EXECUTABLE}" build -trimpath -o "${NEXTCMD_NATIVE_OUTPUT}" .
    DEPENDS ${NEXTCMD_GO_SOURCES}
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Building ${NEXTCMD_NATIVE_OUTPUT}"
    VERBATIM
)
add_custom_target(nextcmd ALL DEPENDS "${NEXTCMD_NATIVE_OUTPUT}")

# This explicit target always creates a directly runnable executable in the
# repository root, regardless of the configured default output location.
set(NEXTCMD_ROOT_OUTPUT "${CMAKE_SOURCE_DIR}/${NEXTCMD_EXECUTABLE_NAME}")
if(NEXTCMD_OUTPUT_IN_SOURCE_ROOT)
    add_custom_target(nextcmd-root DEPENDS "${NEXTCMD_NATIVE_OUTPUT}")
else()
    add_custom_command(
        OUTPUT "${NEXTCMD_ROOT_OUTPUT}"
        COMMAND "${CMAKE_COMMAND}" -E env
            "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
            "${GO_EXECUTABLE}" build -trimpath -o "${NEXTCMD_ROOT_OUTPUT}" .
        DEPENDS ${NEXTCMD_GO_SOURCES}
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        COMMENT "Building ${NEXTCMD_ROOT_OUTPUT}"
        VERBATIM
    )
    add_custom_target(nextcmd-root DEPENDS "${NEXTCMD_ROOT_OUTPUT}")
endif()
set_property(DIRECTORY APPEND PROPERTY ADDITIONAL_CLEAN_FILES "${NEXTCMD_ROOT_OUTPUT}")

add_custom_target(format
    COMMAND "${GO_EXECUTABLE}" fmt ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Formatting Go packages"
    VERBATIM
)

add_custom_target(vet
    COMMAND "${CMAKE_COMMAND}" -E env
        "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
        "${GO_EXECUTABLE}" vet ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running go vet"
    VERBATIM
)

add_custom_target(test
    COMMAND "${CMAKE_COMMAND}" -E env
        "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
        "${GO_EXECUTABLE}" test ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running unit tests"
    VERBATIM
)

add_custom_target(test-race
    COMMAND "${CMAKE_COMMAND}" -E env
        "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
        "${GO_EXECUTABLE}" test -race ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running tests with the race detector (requires CGO)"
    VERBATIM
)

function(add_nextcmd_cross_build target_name target_os target_arch executable_name)
    add_custom_target(${target_name}
        COMMAND "${CMAKE_COMMAND}" -E make_directory "${NEXTCMD_ARTIFACT_DIR}"
        COMMAND "${CMAKE_COMMAND}" -E env
            "GOCACHE=${NEXTCMD_GO_CACHE_DIR}"
            "GOOS=${target_os}"
            "GOARCH=${target_arch}"
            "CGO_ENABLED=0"
            "${GO_EXECUTABLE}" build -trimpath -o "${NEXTCMD_ARTIFACT_DIR}/${executable_name}" .
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        COMMENT "Cross-building NextCmd for ${target_os}/${target_arch}"
        VERBATIM
    )
endfunction()

add_nextcmd_cross_build(build-windows-amd64 windows amd64 NextCmd-windows-amd64.exe)
add_nextcmd_cross_build(build-windows-arm64 windows arm64 NextCmd-windows-arm64.exe)
add_nextcmd_cross_build(build-linux-amd64 linux amd64 NextCmd-linux-amd64)
add_nextcmd_cross_build(build-linux-arm64 linux arm64 NextCmd-linux-arm64)
add_nextcmd_cross_build(build-darwin-amd64 darwin amd64 NextCmd-darwin-amd64)
add_nextcmd_cross_build(build-darwin-arm64 darwin arm64 NextCmd-darwin-arm64)

add_custom_target(build-all-platforms
    DEPENDS
        build-windows-amd64 build-windows-arm64
        build-linux-amd64 build-linux-arm64
        build-darwin-amd64 build-darwin-arm64
)
