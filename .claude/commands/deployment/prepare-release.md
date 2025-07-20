# Prepare Release Command

<instructions>
  <context>
    Comprehensive release preparation with automated quality checks, version management, changelog generation, and deployment readiness validation. Ensures consistent, safe, and documented releases.
  </context>
  
  <requirements>
    - Version control system with tagging capability
    - Automated testing and quality gate pipelines
    - Changelog and release documentation tools
    - Deployment infrastructure and tooling
  </requirements>
  
  <execution>
    1. **Quality Gate Validation**
       - Execute full test suite (unit, integration, e2e)
       - Run security scanning and vulnerability checks
       - Validate code quality and coverage metrics
       - Verify performance benchmarks
    
    2. **Version Management**
       - Analyze changes for semantic versioning
       - Update version numbers in all relevant files
       - Generate version tags and release branches
       - Update dependency versions if needed
    
    3. **Changelog Generation**
       - Collect commits since last release
       - Categorize changes (features, fixes, breaking)
       - Generate human-readable changelog
       - Create release notes and migration guides
    
    4. **Build and Artifact Creation**
       - Execute production build process
       - Generate deployment artifacts
       - Create container images if applicable
       - Sign and verify build artifacts
    
    5. **Documentation Updates**
       - Update API documentation
       - Refresh user guides and tutorials
       - Update installation and setup instructions
       - Generate deployment documentation
    
    6. **Release Validation**
       - Validate deployment readiness
       - Verify rollback procedures
       - Check monitoring and alerting setup
       - Confirm backup and recovery plans
  </execution>
  
  <validation>
    - [ ] All quality gates passed
    - [ ] Version numbers updated
    - [ ] Changelog generated
    - [ ] Build artifacts created
    - [ ] Documentation updated
    - [ ] Security scan clean
    - [ ] Performance benchmarks met
    - [ ] Deployment plan validated
  </validation>
  
  <examples>
    ```bash
    # Prepare patch release
    /deploy:prepare-release --type=patch
    
    # Release preparation:
    # - Version bump: 1.2.3 → 1.2.4
    # - Bug fix changelog
    # - Security updates
    # - Quick deployment
    ```
    
    ```bash
    # Prepare major release
    /deploy:prepare-release --type=major
    
    # Major release preparation:
    # - Version bump: 1.2.3 → 2.0.0
    # - Breaking changes documentation
    # - Migration guides
    # - Extensive testing
    ```
    
    ```bash
    # Prepare release with custom notes
    /deploy:prepare-release --notes="Critical security update"
    
    # Custom release preparation:
    # - Special release notes
    # - Targeted deployment
    # - Emergency procedures
    # - Stakeholder communication
    ```
  </examples>
</instructions>