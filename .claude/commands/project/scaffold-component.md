# Scaffold Component Command

<instructions>
  <context>
    Rapid component creation with boilerplate code, tests, and documentation. Analyzes existing component patterns and generates consistent, production-ready components following project conventions.
  </context>
  
  <requirements>
    - Component library or framework configured
    - Style system in place (CSS, SCSS, styled-components, etc.)
    - Testing framework configured
    - Storybook or similar documentation tool (optional)
  </requirements>
  
  <execution>
    1. **Pattern Analysis**
       - Scan existing components for patterns
       - Identify naming conventions
       - Analyze prop patterns and TypeScript interfaces
       - Detect styling approach and file organization
    
    2. **Component Generation**
       - Create component file with proper structure
       - Generate TypeScript interfaces/PropTypes
       - Add default props and proper typing
       - Include accessibility attributes
    
    3. **Styling**
       - Create corresponding style files
       - Add responsive design considerations
       - Include theme variables and design tokens
       - Follow BEM or project naming conventions
    
    4. **Testing**
       - Generate unit tests with common scenarios
       - Add accessibility testing
       - Create visual regression tests if configured
       - Include interaction testing
    
    5. **Documentation**
       - Generate Storybook stories
       - Add JSDoc or similar documentation
       - Create usage examples
       - Document props and variants
    
    6. **Integration**
       - Update component index files
       - Add to component library exports
       - Update type definitions
       - Add to design system documentation
  </execution>
  
  <validation>
    - [ ] Component renders without errors
    - [ ] All props work as expected
    - [ ] TypeScript compilation passes
    - [ ] Tests cover main use cases
    - [ ] Accessibility standards met
    - [ ] Responsive design works
    - [ ] Storybook documentation complete
    - [ ] Follows project conventions
  </validation>
  
  <examples>
    ```bash
    # Basic component scaffolding
    /project:scaffold-component Button
    
    # Generated files:
    # - components/Button/Button.tsx
    # - components/Button/Button.styles.ts
    # - components/Button/Button.test.tsx
    # - components/Button/Button.stories.tsx
    # - components/Button/index.ts
    ```
    
    ```bash
    # Complex component with variants
    /project:scaffold-component Card --variants=default,outlined,elevated
    
    # Generated with:
    # - Multiple style variants
    # - Comprehensive prop types
    # - Stories for each variant
    # - Accessibility considerations
    ```
    
    ```bash
    # Form component with validation
    /project:scaffold-component ContactForm --type=form --validation=yup
    
    # Generated with:
    # - Form state management
    # - Validation schema
    # - Error handling
    # - Submission logic
    ```
  </examples>
</instructions>