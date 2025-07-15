# Claude Pilot Theming Standards

## Overview

This document establishes comprehensive theming standards for Claude Pilot, ensuring consistent visual identity across CLI and TUI interfaces. The theme is centered around Claude orange (#FF6B35) as the primary brand color, with a carefully crafted palette designed for accessibility and visual hierarchy.

## Color Palette

### Primary Brand Colors

- **ClaudePrimary** (`#FF6B35`) - Main brand color, used for primary actions, highlights, and brand elements
- **ClaudePrimaryLight** (`#FF8A65`) - Lighter orange for hover states and subtle highlights
- **ClaudePrimaryDark** (`#E55A2B`) - Darker orange for pressed states and emphasis

### Secondary Colors

- **ClaudeSecondary** (`#6BB6FF`) - Blue accent for links, secondary actions, and complementary elements
- **ClaudeSecondaryLight** (`#8FC7FF`) - Light blue for hover states
- **ClaudeSecondaryDark** (`#4A90E2`) - Darker blue for pressed states

### Neutral Colors

#### Backgrounds

- **BackgroundPrimary** (`#2C3E50`) - Main dark background
- **BackgroundSecondary** (`#34495E`) - Cards, panels, elevated surfaces
- **BackgroundSurface** (`#4A5568`) - Subtle elevation, borders
- **BackgroundSurfaceLight** (`#718096`) - Light surfaces, disabled states

#### Text Colors

- **TextPrimary** (`#FFFFFF`) - High contrast text, headers, important content
- **TextSecondary** (`#D5DBDB`) - Medium contrast text, body content
- **TextMuted** (`#AEB6BF`) - Low contrast text, captions, metadata
- **TextDim** (`#85929E`) - Very low contrast, timestamps, subtle info
- **TextDisabled** (`#718096`) - Disabled text and elements

### Semantic Colors

#### Status Indicators

- **SuccessColor** (`#2ECC71`) - Success states, active sessions, completed actions
- **WarningColor** (`#F39C12`) - Warnings, inactive states, pending actions
- **ErrorColor** (`#E74C3C`) - Errors, failed states, dangerous actions
- **InfoColor** (`#5DADE2`) - Information, neutral states, help text

#### Interactive States

- **HoverColor** - ClaudePrimaryLight for hover effects
- **FocusColor** - ClaudePrimary for focus indicators
- **ActiveColor** - ClaudePrimaryDark for active/pressed states
- **SelectedColor** - ClaudePrimary for selected items
- **DisabledColor** (`#4A5568`) - Disabled elements

## Component Styling Standards

### Buttons

#### Primary Button

```go
ButtonPrimaryStyle = lipgloss.NewStyle().
    Foreground(TextPrimary).
    Background(ActionPrimary).
    Bold(true).
    Padding(0, 2).
    Margin(0, 1)
```

#### Button States

- **Focused**: Inverted colors with focus border
- **Disabled**: Muted colors, no interaction
- **Hover**: Lighter background color

### Forms

#### Input Fields

- **Default**: Normal border with muted color
- **Focused**: Border changes to FocusColor
- **Error**: Border changes to ErrorColor
- **Disabled**: Muted appearance

#### Labels

- **Default**: Primary text, bold
- **Required**: Same as default (use * indicator)
- **Error**: ErrorColor for validation messages

### Tables

#### Structure

- **Headers**: Primary background with white text, centered
- **Cells**: Left-aligned, secondary text
- **Selected Rows**: Primary background with white text
- **Alternate Rows**: Subtle background for better readability

#### Status Cells

Use semantic colors for status indicators:

- Success states: SuccessColor
- Warning states: WarningColor
- Error states: ErrorColor
- Info states: InfoColor

### Cards and Containers

#### Basic Cards

- Rounded border with surface color
- Secondary background
- Proper padding and margins

#### Semantic Containers

- **InfoBox**: Info color border, normal border style
- **SuccessBox**: Success color border, rounded style
- **WarningBox**: Warning color border, thick style
- **ErrorBox**: Error color border, double style

## Accessibility Guidelines

### Contrast Requirements

All color combinations must meet WCAG 2.1 AA standards:

- Normal text: 4.5:1 contrast ratio minimum
- Large text: 3:1 contrast ratio minimum
- UI components: 3:1 contrast ratio minimum

### Color-Blind Considerations

- Never rely solely on color to convey information
- Use icons, text labels, or patterns alongside colors
- Test with color-blind simulation tools
- Ensure sufficient contrast between different states

### Terminal Compatibility

- Test across different terminal emulators
- Provide fallbacks for terminals with limited color support
- Ensure readability in both dark and light terminal themes

## Usage Guidelines

### When to Use Primary Colors

- **ClaudePrimary**: Brand elements, primary actions, selected states, active navigation
- **ClaudeSecondary**: Links, secondary actions, complementary highlights

### Interactive State Hierarchy

1. **Default**: Base colors as defined
2. **Hover**: Lighter variants of base colors
3. **Focus**: Primary color with additional indicators (borders, underlines)
4. **Active/Pressed**: Darker variants of base colors
5. **Disabled**: Muted gray colors

### Typography Hierarchy

1. **Titles**: ClaudePrimary, bold, larger size
2. **Headers**: TextPrimary, bold, underlined
3. **Body Text**: TextSecondary, normal weight
4. **Captions**: TextMuted, smaller size
5. **Metadata**: TextDim, smallest size

## Responsive Behavior

### Breakpoints

- **Small** (`< 80 chars`): Compact layout, minimal padding
- **Medium** (`80-120 chars`): Balanced layout, normal padding
- **Large** (`> 120 chars`): Full layout, generous padding

### Adaptive Styling

Use responsive utilities:

- `AdaptiveWidth()` for component sizing
- `ResponsivePadding()` for spacing
- `GetTableColumnWidths()` for table layouts

## Migration Guide

### From Legacy Styling

1. Replace direct `fatih/color` usage with `lipgloss` styles
2. Use semantic color constants instead of hardcoded values
3. Apply consistent interactive states
4. Implement responsive behavior

### Best Practices

1. Always use theme constants, never hardcode colors
2. Apply semantic meaning to color choices
3. Test across different terminal sizes
4. Maintain consistency between CLI and TUI
5. Document any custom styling decisions

## Code Examples

### Basic Usage

```go
// Good: Using theme constants
title := styles.TitleStyle.Render("Claude Pilot")
success := styles.Success("Operation completed")

// Bad: Hardcoded colors
title := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B35")).Render("Title")
```

### Responsive Components

```go
// Adaptive button based on terminal width
func CreateButton(text string, width int) string {
    style := styles.ButtonPrimaryStyle
    if width < styles.BreakpointSmall {
        style = style.Padding(0, 1) // Reduce padding for small screens
    }
    return style.Render(text)
}
```

### Status Indicators Functions

```go
// Using contextual colors
func FormatStatus(status string) string {
    color := styles.GetContextualColor("status", status)
    return lipgloss.NewStyle().Foreground(color).Bold(true).Render(status)
}
```

## Validation Checklist

- [ ] All colors meet accessibility contrast requirements
- [ ] Interactive states are clearly defined and consistent
- [ ] Components work across all breakpoints
- [ ] No hardcoded colors in implementation
- [ ] Semantic meaning is applied consistently
- [ ] Terminal compatibility is verified
- [ ] Brand identity is maintained throughout
