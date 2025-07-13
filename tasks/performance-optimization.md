# Claude Pilot Performance Optimization Plan

## Overview

Optimize Claude Pilot CLI performance by eliminating unnecessary abstractions and improving session management efficiency.

## Phase 1: Architecture Simplification (High Priority)

### Task 1: Eliminate SessionManager Layer

- **Problem**: SessionManager adds unnecessary indirection
- **Solution**: Move SessionManager logic directly into commands
- **Files**: `internal/manager/session_manager.go`, `cmd/*.go`
- **Impact**: Reduces function call overhead and complexity

### Task 2: Simplify Factory Pattern

- **Problem**: MultiplexerFactory is over-engineered for 2 backends
- **Solution**: Replace with simple detection function
- **Files**: `internal/multiplexer/factory.go`, related imports
- **Impact**: Reduces object allocation and call overhead

### Task 3: Consolidate Interfaces

- **Problem**: Multiple small interfaces create abstraction overhead
- **Solution**: Merge related interfaces where beneficial
- **Files**: `internal/interfaces/*.go`
- **Impact**: Reduces interface dispatch overhead

## Phase 2: Performance Optimizations (Medium Priority)

### Task 4: Add Session Name Indexing

- **Problem**: Linear search through sessions by name
- **Solution**: Add name-to-ID mapping in repository
- **Files**: `internal/storage/file_repository.go`
- **Impact**: O(1) name lookups instead of O(n)

### Task 5: Implement Multiplexer Caching

- **Problem**: Repeated multiplexer instance creation
- **Solution**: Cache multiplexer instances per backend
- **Files**: Multiplexer creation sites
- **Impact**: Reduces object allocation

### Task 6: Batch Status Updates

- **Problem**: Individual status checks for each session
- **Solution**: Batch status updates where possible
- **Files**: Session listing and status operations
- **Impact**: Reduces system calls

## Phase 3: Micro-optimizations (Low Priority)

### Task 7: Fix Session Creation I/O

- **Problem**: Inefficient file operations during session creation
- **Solution**: Optimize file writing and directory operations
- **Files**: `internal/storage/file_repository.go`
- **Impact**: Faster session creation

### Task 8: Pre-allocate Slices

- **Problem**: Slice growth during session operations
- **Solution**: Pre-allocate slices with known capacity
- **Files**: Session listing and filtering operations
- **Impact**: Reduces memory allocations

## Success Criteria

- All phases maintain backward compatibility
- Build passes after each task completion
- No functional regressions
- Measurable performance improvement in session operations
