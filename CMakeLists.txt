cmake_minimum_required(VERSION 3.20)

project(NextCmd LANGUAGES NONE)

find_program(GO_EXECUTABLE NAMES go REQUIRED)

set(NEXTCMD_OUTPUT_DIR "${CMAKE_BINARY_DIR}/bin")
if(WIN32)
    set(NEXTCMD_EXECUTABLE "nextcmd.exe")
else()
    set(NEXTCMD_EXECUTABLE "nextcmd")
endif()

# Go remains the native build system; CMake provides a familiar, portable facade.
add_custom_target(nextcmd ALL
    COMMAND "${CMAKE_COMMAND}" -E make_directory "${NEXTCMD_OUTPUT_DIR}"
    COMMAND "${GO_EXECUTABLE}" build -o "${NEXTCMD_OUTPUT_DIR}/${NEXTCMD_EXECUTABLE}" .
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Building NextCmd"
    VERBATIM
)

add_custom_target(format
    COMMAND "${GO_EXECUTABLE}" fmt ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Formatting Go packages"
    VERBATIM
)

add_custom_target(vet
    COMMAND "${GO_EXECUTABLE}" vet ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running go vet"
    VERBATIM
)

add_custom_target(test
    COMMAND "${GO_EXECUTABLE}" test ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running unit tests"
    VERBATIM
)

add_custom_target(test-race
    COMMAND "${GO_EXECUTABLE}" test -race ./...
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    COMMENT "Running tests with the race detector (requires CGO)"
    VERBATIM
)

function(add_nextcmd_cross_build target_name target_os target_arch executable_name)
    add_custom_target(${target_name}
        COMMAND "${CMAKE_COMMAND}" -E make_directory "${NEXTCMD_OUTPUT_DIR}"
        COMMAND "${CMAKE_COMMAND}" -E env
            "GOOS=${target_os}"
            "GOARCH=${target_arch}"
            "CGO_ENABLED=0"
            "${GO_EXECUTABLE}" build -o "${NEXTCMD_OUTPUT_DIR}/${executable_name}" .
        WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
        COMMENT "Cross-building NextCmd for ${target_os}/${target_arch}"
        VERBATIM
    )
endfunction()

add_nextcmd_cross_build(build-windows-amd64 windows amd64 nextcmd-windows-amd64.exe)
add_nextcmd_cross_build(build-windows-arm64 windows arm64 nextcmd-windows-arm64.exe)
add_nextcmd_cross_build(build-linux-amd64 linux amd64 nextcmd-linux-amd64)
add_nextcmd_cross_build(build-linux-arm64 linux arm64 nextcmd-linux-arm64)
add_nextcmd_cross_build(build-darwin-amd64 darwin amd64 nextcmd-darwin-amd64)
add_nextcmd_cross_build(build-darwin-arm64 darwin arm64 nextcmd-darwin-arm64)

add_custom_target(build-all-platforms
    DEPENDS
        build-windows-amd64 build-windows-arm64
        build-linux-amd64 build-linux-arm64
        build-darwin-amd64 build-darwin-arm64
)
