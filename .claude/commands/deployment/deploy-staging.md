# Deploy Staging Command

<instructions>
  <context>
    Automated staging deployment with comprehensive validation, smoke testing, and rollback capabilities. Provides safe deployment testing environment before production release.
  </context>
  
  <requirements>
    - Staging environment configured and accessible
    - Deployment pipeline and tools configured
    - Monitoring and logging systems active
    - Rollback procedures documented and tested
  </requirements>
  
  <execution>
    1. **Pre-deployment Validation**
       - Verify staging environment health
       - Check deployment artifacts integrity
       - Validate configuration and secrets
       - Confirm resource availability
    
    2. **Deployment Execution**
       - Execute deployment pipeline
       - Monitor deployment progress
       - Handle deployment failures gracefully
       - Maintain deployment logs and audit trail
    
    3. **Service Validation**
       - Verify service startup and health checks
       - Test API endpoints and functionality
       - Validate database connections and migrations
       - Check external service integrations
    
    4. **Smoke Testing**
       - Execute critical user journeys
       - Validate key business functionality
       - Test authentication and authorization
       - Verify performance characteristics
    
    5. **Monitoring Setup**
       - Configure environment-specific monitoring
       - Set up alerting for staging issues
       - Enable logging and tracing
       - Create deployment dashboards
    
    6. **Validation and Sign-off**
       - Generate deployment report
       - Document any issues or concerns
       - Collect stakeholder approval
       - Prepare production deployment plan
  </execution>
  
  <validation>
    - [ ] Deployment completed successfully
    - [ ] All services healthy and responsive
    - [ ] Smoke tests passed
    - [ ] Database migrations successful
    - [ ] External integrations working
    - [ ] Monitoring and alerts active
    - [ ] Performance metrics acceptable
    - [ ] Rollback procedure verified
  </validation>
  
  <examples>
    ```bash
    # Deploy current branch to staging
    /deploy:deploy-staging
    
    # Deployment process:
    # - Environment validation
    # - Artifact deployment
    # - Service health checks
    # - Smoke test execution
    ```
    
    ```bash
    # Deploy specific version to staging
    /deploy:deploy-staging --version=1.2.3
    
    # Version-specific deployment:
    # - Tagged release deployment
    # - Version validation
    # - Rollback planning
    # - Change documentation
    ```
    
    ```bash
    # Deploy with custom configuration
    /deploy:deploy-staging --config=feature-flags-enabled
    
    # Custom deployment:
    # - Feature flag configuration
    # - A/B testing setup
    # - Experimental features
    # - Controlled rollout
    ```
  </examples>
</instructions>