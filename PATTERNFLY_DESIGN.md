# PatternFly Design System Implementation

This document outlines how the OVN-Kubernetes Trace Viewer follows PatternFly and Red Hat UX design principles.

## Design System References

- **PatternFly**: https://www.patternfly.org/
- **Red Hat UX**: https://ux.redhat.com/
- **PatternFly React**: https://www.patternfly.org/get-started/develop#react

## Component Usage

### Layout Components

#### Page Structure
```tsx
<Page header={masthead}>
  <Drawer isExpanded={true} isInline>
    <DrawerContent panelContent={configPanel}>
      <PageSection>
        {/* Main content */}
      </PageSection>
    </DrawerContent>
  </Drawer>
</Page>
```

**Why**: PatternFly's Page component provides consistent header, navigation, and content structure. The Drawer component enables a collapsible configuration panel without losing context.

**Red Hat UX Principle**: *Progressive Disclosure* - Show configuration when needed, hide when focusing on visualization.

#### Masthead (Header)
```tsx
<Masthead>
  <MastheadMain>
    <MastheadBrand>
      <Brand src={logo} alt="OVN Logo">
        <Title>OVN-Kubernetes Trace Viewer</Title>
      </Brand>
    </MastheadBrand>
  </MastheadMain>
  <MastheadContent>
    {/* Description */}
  </MastheadContent>
</Masthead>
```

**Why**: Consistent branding and navigation across Red Hat products.

**Red Hat UX Principle**: *Consistency* - Users familiar with OpenShift will recognize the layout.

### Form Components

#### Trace Configuration Form
```tsx
<Form>
  <FormSection title="Network Endpoints">
    <FormGroup label="Source" isRequired>
      <TextInput {...props} />
      <FormHelperText>
        <HelperText>
          <HelperTextItem>Pod name, IP, or namespace/pod</HelperTextItem>
        </HelperText>
      </FormHelperText>
    </FormGroup>
  </FormSection>
  <Button variant="primary" isBlock icon={<PlayIcon />}>
    Execute Trace
  </Button>
</Form>
```

**Components Used**:
- `Form` - Container for form elements
- `FormSection` - Groups related fields
- `FormGroup` - Individual field with label
- `FormHelperText` - Contextual help
- `Button` - Primary action with icon

**Red Hat UX Principles**:
- *Clarity* - Helper text explains what each field accepts
- *Efficiency* - Grouped sections reduce cognitive load
- *Accessibility* - Proper labels, ARIA attributes

### Visualization Components

#### Flow Topology Display
```tsx
<Card isFullHeight>
  <CardTitle>Network Topology</CardTitle>
  <CardBody>
    <ReactFlow {...props}>
      <Background variant={BackgroundVariant.Dots} />
      <Controls />
    </ReactFlow>
  </CardBody>
</Card>
```

**Why**: Cards provide visual containment and hierarchy. ReactFlow integrates with PatternFly styling.

#### Flow Steps Timeline
```tsx
<Timeline>
  {steps.map(step => (
    <TimelineItem>
      <Text component={TextVariants.h4}>
        Step {step.number}: <Badge>{step.action}</Badge>
      </Text>
      <Text>{step.description}</Text>
    </TimelineItem>
  ))}
</Timeline>
```

**Components Used**:
- `Timeline` - Sequential step visualization
- `Badge` - Highlight key information (action, protocol)
- `Text` - Typography components with semantic variants

**Red Hat UX Principle**: *Understandability* - Timeline shows clear progression of packet flow.

### Status & Feedback

#### Trace Status
```tsx
<Label icon={<CheckCircleIcon />} color="green">
  SUCCESS
</Label>
```

**Status Indicators**:
- ✅ Success → Green label with CheckCircleIcon
- ❌ Dropped → Red label with TimesCircleIcon  
- ⚠️ Error → Orange label with ExclamationCircleIcon

**Red Hat UX Principle**: *Feedback* - Immediate visual indication of trace result.

#### Loading States
```tsx
<EmptyState>
  <EmptyStateHeader
    icon={<EmptyStateIcon icon={Spinner} />}
    titleText="Executing trace..."
  />
  <EmptyStateBody>
    Please wait while we trace the network path...
  </EmptyStateBody>
</EmptyState>
```

**Red Hat UX Principle**: *Responsiveness* - Always show system state, never leave user wondering.

#### Error Handling
```tsx
<Alert
  variant={AlertVariant.danger}
  title="Trace Execution Error"
  actionClose={<AlertActionCloseButton onClose={handleClose} />}
>
  {errorMessage}
</Alert>
```

**Red Hat UX Principle**: *Error Prevention & Recovery* - Clear error messages with dismissal action.

## Color Palette

Following PatternFly's color system:

### Node Type Colors
```css
--switch-color: #06c;      /* PatternFly blue-400 */
--router-color: #a30000;   /* PatternFly red-200 */
--port-color: #8476d1;     /* PatternFly purple-300 */
--acl-color: #f0ab00;      /* PatternFly gold-400 */
--default-color: #6a6e73;  /* PatternFly grey-600 */
```

### Status Colors
```css
--success: var(--pf-v5-global--success-color--100);  /* Green */
--danger: var(--pf-v5-global--danger-color--100);    /* Red */
--warning: var(--pf-v5-global--warning-color--100);  /* Orange */
```

**Why**: Using PatternFly color tokens ensures:
- Accessibility (WCAG AA contrast ratios)
- Theme consistency
- Dark mode support (future)

## Typography

### Font Stack
```css
font-family: "RedHatDisplay", "Overpass", overpass, helvetica, arial, sans-serif;
```

**Red Hat Display** - Primary font for headings and UI
**Overpass** - Fallback sans-serif

### Text Components
```tsx
<Text component={TextVariants.h1}>Main Title</Text>
<Text component={TextVariants.p}>Body text</Text>
<Text component={TextVariants.small}>Helper text</Text>
```

**Why**: Semantic HTML with consistent sizing and spacing.

## Spacing & Layout

PatternFly uses an 8px grid system:
- Spacing tokens: `--pf-v5-global--spacer--{xs|sm|md|lg|xl}`
- Components automatically apply consistent spacing

### Grid Layout
```tsx
<Split hasGutter>
  <SplitItem>Source</SplitItem>
  <SplitItem>→</SplitItem>
  <SplitItem>Destination</SplitItem>
</Split>
```

## Responsive Design

### Breakpoints
- **Mobile**: < 768px - Single column, stacked panels
- **Tablet**: 768px - 1200px - Responsive drawer
- **Desktop**: > 1200px - Full layout with side panel

```css
@media (max-width: 1200px) {
  .flow-content {
    grid-template-columns: 1fr;
    grid-template-rows: 1fr auto;
  }
}
```

## Accessibility

### WCAG 2.1 AA Compliance

✅ **Keyboard Navigation**
- All interactive elements focusable
- Logical tab order
- Drawer toggle with keyboard

✅ **Screen Reader Support**
- Semantic HTML (headings, labels, descriptions)
- ARIA labels on icons
- Form field associations

✅ **Color Contrast**
- All text meets 4.5:1 ratio
- Interactive elements meet 3:1 ratio

✅ **Focus Indicators**
- PatternFly provides visible focus rings

### Testing
```bash
# Run accessibility audit
npm run build
npx pa11y-ci http://localhost:3000
```

## Interaction Patterns

### Primary Action Flow
1. **User opens app** → Empty state with clear instructions
2. **User configures trace** → Drawer panel with form
3. **User submits** → Loading spinner with message
4. **Trace completes** → Results with status indicator
5. **User explores** → Interactive graph + timeline

### Progressive Disclosure
- Configuration panel can be collapsed
- Focus on visualization when not configuring
- Toggle button to re-open panel

### Feedback Loop
- Loading states during async operations
- Success/error alerts with clear messaging
- Status badges throughout UI

## Component Inventory

### Layout
- ✅ Page
- ✅ Masthead
- ✅ Drawer
- ✅ Toolbar
- ✅ Card

### Forms
- ✅ Form / FormGroup / FormSection
- ✅ TextInput
- ✅ FormSelect
- ✅ Button
- ✅ FormHelperText

### Data Display
- ✅ Timeline
- ✅ DescriptionList
- ✅ Badge
- ✅ Label

### Feedback
- ✅ Alert
- ✅ EmptyState
- ✅ Spinner

### Typography
- ✅ Title
- ✅ Text / TextContent

### Icons
- ✅ @patternfly/react-icons (700+ icons)

## Design Principles Applied

### 1. Consistency
**Applied**: Using PatternFly components ensures visual and behavioral consistency with OpenShift console and other Red Hat products.

**Example**: Masthead layout, form patterns, button styles all match OpenShift UX.

### 2. Clarity
**Applied**: Clear labeling, helper text, and status indicators.

**Example**: Each form field has a label and helper text explaining what to enter.

### 3. Efficiency
**Applied**: Minimal clicks, logical grouping, keyboard shortcuts.

**Example**: Form sections group related fields. Primary action is a single button.

### 4. Feedback
**Applied**: Loading states, success/error alerts, status badges.

**Example**: Spinner during trace execution, Alert on error, Label showing trace status.

### 5. Accessibility
**Applied**: Semantic HTML, ARIA labels, keyboard navigation, color contrast.

**Example**: All form fields have proper labels, focus indicators, screen reader support.

### 6. Progressive Disclosure
**Applied**: Collapsible drawer, empty states, timeline expansion.

**Example**: Configuration panel can be hidden to focus on visualization.

## Future Enhancements

### Phase 2 Additions
- [ ] Data table for trace history
- [ ] Pagination for large datasets
- [ ] Advanced filters
- [ ] Export functionality (PatternFly button group)

### Phase 3 Additions
- [ ] Tabs for multiple traces
- [ ] Modal for trace comparison
- [ ] Notification drawer for background traces

### Phase 4 - OpenShift Console Plugin
- [ ] Full console integration
- [ ] Navigation items
- [ ] Action menus on resource pages
- [ ] Dark mode support

## Resources

### PatternFly
- [Get Started](https://www.patternfly.org/get-started/develop)
- [Components](https://www.patternfly.org/components/all-components)
- [Design Guidelines](https://www.patternfly.org/design-guidelines/styles/colors)
- [Accessibility](https://www.patternfly.org/accessibility/accessibility-fundamentals)

### Red Hat UX
- [Design System](https://ux.redhat.com/)
- [Foundations](https://ux.redhat.com/foundations/)
- [Patterns](https://ux.redhat.com/patterns/)

### Code Examples
- [PatternFly React Seed](https://github.com/patternfly/patternfly-react-seed)
- [OpenShift Console](https://github.com/openshift/console)

---

**Updated**: October 26, 2026  
**PatternFly Version**: 5.4.1  
**Compliant With**: Red Hat UX Design System
