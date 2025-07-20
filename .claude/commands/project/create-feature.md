# Create Feature Command

<instructions>
  <context>
    Full-stack feature development with comprehensive scaffolding, testing, and documentation. Creates feature branch, implements core functionality, adds tests, and updates documentation in a single workflow.
  </context>
  
  <requirements>
    - Git repository initialized
    - Development environment configured
    - Package manager available (npm, yarn, pip, go mod, etc.)
    - Test framework configured
    - Documentation structure exists
  </requirements>
  
  <execution>
    1. **Feature Planning**
       - Analyze existing codebase structure
       - Identify integration points and dependencies
       - Create feature specification document
       - Estimate complexity and timeline
    
    2. **Branch Management**
       - Create feature branch: `feat/[feature-name]`
       - Set up branch protection if needed
       - Link to issue tracker (Linear, GitHub Issues)
    
    3. **Core Implementation**
       - Create main feature files following project conventions
       - Implement core functionality with proper error handling
       - Add input validation and security checks
       - Follow established patterns and style guides
    
    4. **Test Implementation**
       - Generate unit tests for all new functions
       - Create integration tests for feature workflows
       - Add end-to-end tests for user-facing features
       - Ensure minimum coverage thresholds are met
    
    5. **Documentation**
       - Update API documentation if applicable
       - Add inline code documentation
       - Create user-facing documentation
       - Update changelog and migration guides
    
    6. **Quality Assurance**
       - Run linting and formatting tools
       - Execute full test suite
       - Perform security scanning
       - Validate performance requirements
  </execution>
  
  <validation>
    - [ ] Feature works as specified
    - [ ] All tests pass (unit, integration, e2e)
    - [ ] Code coverage meets project standards
    - [ ] Documentation is complete and accurate
    - [ ] Security scan passes
    - [ ] Performance benchmarks met
    - [ ] Code review checklist complete
    - [ ] No breaking changes introduced
  </validation>
  
  <examples>
    ```bash
    # Usage example
    /project:create-feature user-authentication
    
    # Expected structure created:
    # - feat/user-authentication branch
    # - src/auth/authentication.{js,ts,py,go}
    # - tests/auth/authentication.test.{js,ts,py,go}
    # - docs/auth/authentication.md
    # - Updated API documentation
    ```
    
    ```bash
    # With additional options
    /project:create-feature payment-processing --type=service --framework=express
    
    # Creates service-specific structure:
    # - services/payment/
    # - middleware/payment-validation.js
    # - routes/payment.js
    # - tests/payment/
    ```
  </examples>
</instructions>