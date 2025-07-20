# Coverage Analysis Command

<instructions>
  <context>
    Detailed test coverage analysis with gap identification, quality assessment, and improvement recommendations. Provides actionable insights for achieving comprehensive test coverage while maintaining test quality.
  </context>
  
  <requirements>
    - Code coverage tools configured (Istanbul, Coverage.py, etc.)
    - Test suite executable and passing
    - Access to source code and test files
    - Coverage reporting tools available
  </requirements>
  
  <execution>
    1. **Coverage Measurement**
       - Execute full test suite with coverage tracking
       - Generate line, branch, and function coverage reports
       - Identify uncovered code paths
       - Calculate coverage percentages by module
    
    2. **Gap Analysis**
       - Identify critical uncovered code sections
       - Analyze missing edge case coverage
       - Find untested error handling paths
       - Detect dead or unreachable code
    
    3. **Quality Assessment**
       - Evaluate test quality vs coverage percentage
       - Identify superficial or ineffective tests
       - Analyze test maintainability
       - Review test performance and execution time
    
    4. **Risk Assessment**
       - Prioritize coverage gaps by business impact
       - Identify high-risk uncovered code
       - Assess coverage trends over time
       - Evaluate coverage stability
    
    5. **Improvement Recommendations**
       - Generate specific test cases for gaps
       - Recommend refactoring for testability
       - Suggest coverage targets by module
       - Provide implementation roadmap
    
    6. **Reporting and Visualization**
       - Create comprehensive coverage reports
       - Generate visual coverage maps
       - Track coverage trends and history
       - Set up coverage monitoring
  </execution>
  
  <validation>
    - [ ] Coverage metrics accurately measured
    - [ ] Critical gaps identified
    - [ ] Test quality assessed
    - [ ] Risk levels assigned
    - [ ] Improvement plan created
    - [ ] Coverage targets set
    - [ ] Monitoring configured
    - [ ] Team reporting established
  </validation>
  
  <examples>
    ```bash
    # Analyze current test coverage
    /test:coverage-analysis
    
    # Coverage report includes:
    # - Overall coverage percentage
    # - Module-by-module breakdown
    # - Critical uncovered lines
    # - Improvement recommendations
    ```
    
    ```bash
    # Analyze specific module coverage
    /test:coverage-analysis --module=auth-service
    
    # Detailed module analysis:
    # - Function-level coverage
    # - Branch coverage analysis
    # - Edge case identification
    # - Specific test recommendations
    ```
    
    ```bash
    # Compare coverage trends
    /test:coverage-analysis --compare=last-week
    
    # Trend analysis:
    # - Coverage changes over time
    # - Regression identification
    # - Improvement tracking
    # - Quality metrics evolution
    ```
  </examples>
</instructions>