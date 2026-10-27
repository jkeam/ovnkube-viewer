# PatternFly Design System

This document covers the PatternFly 5.4 design system implementation for the OVN-Kubernetes Trace Viewer.

## Quick Setup

### Installation

```bash
# Recommended: Use the run script
./run.sh init

# Or manually
cd frontend
npm install
```

This installs:
- `@patternfly/react-core@^5.4.1` - Core components
- `@patternfly/react-icons@^5.4.0` - Icon library
- `@patternfly/react-table@^5.4.1` - Table components
- `@patternfly/patternfly@^5.4.1` - CSS framework

### Running the Application

```bash
# Development mode with hot reload
./run.sh dev

# Or with Docker
./run.sh docker
```

Open http://localhost:3000 to see the PatternFly-styled UI.

## Design System Overview

### Why PatternFly?

1. **OpenShift Alignment** - Phase 4 will be an OpenShift console plugin; using PatternFly now enables component reuse
2. **Enterprise Ready** - WCAG 2.1 AA accessible, professional design
3. **Comprehensive** - 100+ components, 700+ icons, established patterns
4. **Red Hat Standard** - Consistent with OpenShift user experience
5. **Well Maintained** - Active development, regular updates

### Design References

- **PatternFly**: https://www.patternfly.org/
- **Red Hat UX**: https://ux.redhat.com/
- **PatternFly React**: https://www.patternfly.org/get-started/develop#react

## Component Architecture

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

**Red Hat UX Principle**: *Understandability* - Timeline shows clear progression of packet flow.

### Status & Feedback

#### Status Labels
```tsx
<Label icon={<CheckCircleIcon />} color="green">
  SUCCESS
</Label>
```

**Status Indicators**:
- ✅ Success → Green label with CheckCircleIcon
- ❌ Dropped → Red label with TimesCircleIcon  
- ⚠️ Error → Orange label with ExclamationCircleIcon

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

## Visual Design

### Color Palette

Following PatternFly's color system:

#### Node Type Colors
```css
--switch-color: #06c;      /* PatternFly blue-400 */
--router-color: #a30000;   /* PatternFly red-200 */
--port-color: #8476d1;     /* PatternFly purple-300 */
--acl-color: #f0ab00;      /* PatternFly gold-400 */
--default-color: #6a6e73;  /* PatternFly grey-600 */
```

#### Status Colors
```css
--success: var(--pf-v5-global--success-color--100);  /* Green */
--danger: var(--pf-v5-global--danger-color--100);    /* Red */
--warning: var(--pf-v5-global--warning-color--100);  /* Orange */
```

**Why**: Using PatternFly color tokens ensures:
- Accessibility (WCAG AA contrast ratios)
- Theme consistency
- Dark mode support (future)

### Typography

#### Font Stack
```css
font-family: "RedHatDisplay", "Overpass", overpass, helvetica, arial, sans-serif;
```

#### Text Components
```tsx
<Text component={TextVariants.h1}>Main Title</Text>
<Text component={TextVariants.p}>Body text</Text>
<Text component={TextVariants.small}>Helper text</Text>
```

### Spacing & Layout

PatternFly uses an 8px grid system:
- Spacing tokens: `--pf-v5-global--spacer--{xs|sm|md|lg|xl}`
- Components automatically apply consistent spacing

**Spacing scale**: xs (4px), sm (8px), md (16px), lg (24px), xl (32px)

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

### Testing Accessibility

#### Keyboard Navigation
1. Press `Tab` to navigate between elements
2. Press `Enter` or `Space` to activate buttons
3. Press `Esc` to close drawer/modals

#### Screen Reader
1. Enable VoiceOver (Mac): `Cmd + F5`
2. Enable NVDA (Windows): Install from nvda.org
3. Navigate with `Tab` and arrow keys
4. All form fields should announce their labels

#### Automated Testing
```bash
npm run build
npx pa11y-ci http://localhost:3000
```

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

## Customization

### Custom Colors

Override PatternFly colors in `src/App.css`:

```css
/* Custom masthead color */
.pf-v5-c-masthead {
  background: linear-gradient(135deg, #ee0000 0%, #c00 100%);
}

/* Custom node colors */
.react-flow__node.switch {
  background: #06c; /* PatternFly blue */
}
```

### Custom Spacing

```tsx
<div style={{ padding: 'var(--pf-v5-global--spacer--md)' }}>
  {/* 16px padding */}
</div>
```

## Troubleshooting

### PatternFly CSS not loading

**Symptom**: Unstyled components, no red header

**Fix**:
```tsx
// Check src/main.tsx has imports
import '@patternfly/patternfly/patternfly.css'
import '@patternfly/patternfly/patternfly-addons.css'
```

### Components not rendering

**Symptom**: Blank screen, console errors

**Fix**:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Icons missing

**Symptom**: No icons showing

**Fix**:
```tsx
// Check imports
import { PlayIcon, CheckCircleIcon } from '@patternfly/react-icons'
```

## Browser Support

PatternFly 5 supports:
- ✅ Chrome/Edge (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ⚠️ IE 11 (not supported)

## Testing Checklist

- [ ] Header displays with red gradient
- [ ] Logo and title visible
- [ ] Drawer opens and closes
- [ ] Form fields accept input
- [ ] Helper text displays
- [ ] Execute button enabled/disabled based on form validity
- [ ] Loading spinner shows during trace
- [ ] Error alert displays on failure
- [ ] Success status badge shows on success
- [ ] Graph renders in card
- [ ] Timeline shows flow steps
- [ ] Responsive layout works (resize window)
- [ ] Keyboard navigation works (tab through form)
- [ ] Focus indicators visible

## Future Enhancements

### Phase 2 Additions
- [ ] Data table for trace history
- [ ] Pagination for large datasets
- [ ] Advanced filters
- [ ] Export functionality

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
- [Icon Finder](https://www.patternfly.org/design-foundations/icons)

### Red Hat UX
- [Design System](https://ux.redhat.com/)
- [Foundations](https://ux.redhat.com/foundations/)
- [Patterns](https://ux.redhat.com/patterns/)

### Code Examples
- [PatternFly React Seed](https://github.com/patternfly/patternfly-react-seed)
- [OpenShift Console](https://github.com/openshift/console)

### Project Documentation
- [Getting Started Guide](getting-started.md)
- [Implementation Plan](implementation-plan.md)
- [Run Script Guide](run-script-guide.md)

---

**PatternFly Version**: 5.4.1  
**Compliant With**: Red Hat UX Design System, WCAG 2.1 AA
