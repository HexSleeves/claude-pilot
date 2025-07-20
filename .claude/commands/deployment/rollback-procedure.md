# Rollback Procedure Command

<instructions>
  <context>
    Emergency rollback execution with automated recovery, data consistency validation, and incident management. Provides rapid recovery from deployment failures or production issues.
  </context>
  
  <requirements>
    - Rollback infrastructure and automation configured
    - Previous deployment artifacts available
    - Database backup and recovery procedures
    - Incident response team and communication channels
  </requirements>
  
  <execution>
    1. **Incident Assessment**
       - Identify scope and severity of issue
       - Determine rollback necessity and timeline
       - Assess impact on users and business
       - Document incident details and timeline
    
    2. **Rollback Planning**
       - Identify target rollback version
       - Assess data migration requirements
       - Plan service coordination and dependencies
       - Prepare rollback communication
    
    3. **Service Rollback**
       - Stop affected services gracefully
       - Deploy previous version artifacts
       - Restart services in correct order
       - Verify service health and connectivity
    
    4. **Database Rollback**
       - Assess database schema changes
       - Execute database rollback procedures
       - Restore data from backups if needed
       - Verify data consistency and integrity
    
    5. **Validation and Testing**
       - Execute smoke tests on rolled-back system
       - Validate critical functionality
       - Test user authentication and core features
       - Verify external integrations
    
    6. **Communication and Follow-up**
       - Notify stakeholders of rollback completion
       - Document lessons learned
       - Plan incident post-mortem
       - Update rollback procedures based on learnings
  </execution>
  
  <validation>
    - [ ] Rollback completed successfully
    - [ ] All services restored to previous state
    - [ ] Database consistency verified
    - [ ] Smoke tests passed
    - [ ] User impact minimized
    - [ ] Stakeholders notified
    - [ ] Incident documented
    - [ ] Post-mortem scheduled
  </validation>
  
  <examples>
    ```bash
    # Emergency rollback to previous version
    /deploy:rollback-procedure --emergency
    
    # Emergency rollback:
    # - Immediate service rollback
    # - Automated health checks
    # - Stakeholder notification
    # - Incident documentation
    ```
    
    ```bash
    # Rollback to specific version
    /deploy:rollback-procedure --version=1.2.2
    
    # Targeted rollback:
    # - Specific version restoration
    # - Data migration assessment
    # - Comprehensive validation
    # - Controlled recovery
    ```
    
    ```bash
    # Rollback with data preservation
    /deploy:rollback-procedure --preserve-data
    
    # Data-preserving rollback:
    # - Service-only rollback
    # - Database state maintained
    # - Minimal user impact
    # - Gradual recovery
    ```
  </examples>
</instructions>