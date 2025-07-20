# Setup Environment Command

<instructions>
  <context>
    Complete development environment initialization including dependencies, configuration, tooling, and verification. Sets up consistent development environment across team members and CI/CD systems.
  </context>
  
  <requirements>
    - Operating system with package manager
    - Internet connection for downloads
    - Git installed and configured
    - Node.js, Python, Go, or relevant runtime
  </requirements>
  
  <execution>
    1. **Environment Detection**
       - Detect operating system and architecture
       - Check for existing tools and versions
       - Identify package managers available
       - Scan for project type and requirements
    
    2. **Dependency Installation**
       - Install runtime dependencies
       - Set up development dependencies
       - Configure package managers
       - Install global tools and utilities
    
    3. **Configuration Setup**
       - Generate environment configuration files
       - Set up IDE/editor configurations
       - Configure linting and formatting tools
       - Initialize pre-commit hooks
    
    4. **Database and Services**
       - Set up local database instances
       - Configure external service connections
       - Initialize data seeding
       - Set up monitoring and logging
    
    5. **Tool Configuration**
       - Configure test runners
       - Set up build tools and bundlers
       - Initialize documentation tools
       - Configure deployment tools
    
    6. **Verification**
       - Run health checks
       - Execute test suite
       - Verify build process
       - Check service connectivity
  </execution>
  
  <validation>
    - [ ] All dependencies installed successfully
    - [ ] Configuration files generated
    - [ ] Development server starts without errors
    - [ ] Test suite runs and passes
    - [ ] Build process completes successfully
    - [ ] Database connections work
    - [ ] External services accessible
    - [ ] IDE/editor integration working
  </validation>
  
  <examples>
    ```bash
    # Basic setup for new project
    /project:setup-environment
    
    # Actions performed:
    # - npm install / yarn install
    # - Configure .env files
    # - Set up pre-commit hooks
    # - Initialize database
    # - Run initial tests
    ```
    
    ```bash
    # Setup with specific database
    /project:setup-environment --database=postgresql
    
    # Additional actions:
    # - Install PostgreSQL client
    # - Create development database
    # - Run migrations
    # - Seed initial data
    ```
    
    ```bash
    # Team setup with shared configuration
    /project:setup-environment --team-config
    
    # Includes:
    # - Shared IDE settings
    # - Team-wide linting rules
    # - Consistent formatting
    # - Shared Git hooks
    ```
  </examples>
</instructions>