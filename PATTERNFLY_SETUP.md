# PatternFly Setup Guide

Quick guide to get the redesigned UI running with PatternFly.

## Installation

### 1. Install Dependencies

```bash
# Recommended: Use the run script
./run.sh init

# Or manually
cd frontend
npm install
cd ..
```

This installs:
- `@patternfly/react-core@^5.4.1` - Core components
- `@patternfly/react-icons@^5.4.0` - Icon library
- `@patternfly/react-table@^5.4.1` - Table components (future)
- `@patternfly/patternfly@^5.4.1` - CSS framework

### 2. Run the Application

```bash
# Development mode (hot reload)
./run.sh dev

# Or with Docker
./run.sh docker
```

### 3. Access the Application

Open http://localhost:3000 in your browser.

You should see:
- Red Hat branded masthead (header)
- Collapsible drawer with trace configuration form
- PatternFly-styled components throughout

## What to Expect

### First Load
- PatternFly CSS loads (~200KB gzipped)
- Red Hat Display font loads
- Components render with enterprise styling

### UI Elements

#### Header (Masthead)
- Red gradient background (#ee0000)
- OVN logo (SVG shield icon)
- App title: "OVN-Kubernetes Trace Viewer"
- Subtitle: "Visualize network traffic flow..."

#### Configuration Drawer
- Right-side panel (400px default width)
- Collapsible with [x] button
- **Form sections**:
  - Network Endpoints (source, destination)
  - Protocol & Ports (protocol, ports)
- **Primary button**: "Execute Trace" with play icon

#### Main Content Area

**Before trace execution**:
- Empty state with topology icon
- "Ready to trace" message
- Instructions

**During trace**:
- Spinner with "Executing trace..." message
- Loading state

**After trace**:
- Toolbar with status badge and trace metadata
- Card: "Network Topology" with ReactFlow graph
- Timeline: "Flow Steps" with step-by-step breakdown

## Component Preview

### Form Example
```tsx
// PatternFly form components
<FormGroup label="Source" isRequired>
  <TextInput
    placeholder="e.g., pod/web-pod or 10.128.0.5"
    isRequired
  />
  <FormHelperText>
    <HelperTextItem>Pod name, IP, or namespace/pod</HelperTextItem>
  </FormHelperText>
</FormGroup>
```

### Status Badge Example
```tsx
// Success status
<Label icon={<CheckCircleIcon />} color="green">
  SUCCESS
</Label>

// Dropped packet
<Label icon={<TimesCircleIcon />} color="red">
  DROPPED
</Label>
```

### Timeline Example
```tsx
<Timeline>
  <TimelineItem>
    <Text component={TextVariants.h4}>
      Step 1: <Badge>ingress</Badge>
    </Text>
    <Text>Packet enters logical switch...</Text>
  </TimelineItem>
</Timeline>
```

## Customization

### Colors

PatternFly colors are defined in CSS variables. Override in `src/App.css`:

```css
/* Custom masthead color */
.pf-v5-c-masthead {
  background: linear-gradient(135deg, #ee0000 0%, #c00 100%);
}

/* Custom node colors (already set) */
.react-flow__node.switch {
  background: #06c; /* PatternFly blue */
}
```

### Spacing

Use PatternFly spacing tokens:

```tsx
<div style={{ padding: 'var(--pf-v5-global--spacer--md)' }}>
  {/* 16px padding */}
</div>
```

**Spacing scale**: xs (4px), sm (8px), md (16px), lg (24px), xl (32px)

### Typography

Use PatternFly Text components:

```tsx
<Text component={TextVariants.h1}>Main Title</Text>
<Text component={TextVariants.p}>Body text</Text>
<Text component={TextVariants.small}>Helper text</Text>
```

## Responsive Behavior

### Desktop (>1200px)
- Full layout
- Drawer expanded by default
- Graph + timeline side-by-side

### Tablet (768-1200px)
- Drawer collapsible
- Graph above timeline (stacked)

### Mobile (<768px)
- Single column
- Drawer as overlay
- Full-width components

Test by resizing browser window.

## Accessibility Testing

### Keyboard Navigation
1. Press `Tab` to navigate between elements
2. Press `Enter` or `Space` to activate buttons
3. Press `Esc` to close drawer/modals

### Screen Reader
1. Enable VoiceOver (Mac): `Cmd + F5`
2. Enable NVDA (Windows): Install from nvda.org
3. Navigate with `Tab` and arrow keys
4. All form fields should announce their labels

### Color Contrast
All text meets WCAG AA standards:
- Normal text: 4.5:1 contrast ratio
- Large text: 3:1 contrast ratio

Use Chrome DevTools Lighthouse to verify.

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
# Reinstall dependencies
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Font not loading

**Symptom**: Fallback to system font

**Check**: Red Hat Display font should load from PatternFly CSS. No additional setup needed.

### Icons missing

**Symptom**: No icons showing

**Fix**:
```tsx
// Check imports
import { PlayIcon, CheckCircleIcon } from '@patternfly/react-icons'
```

### Drawer not toggling

**Symptom**: Drawer doesn't open/close

**Check**: State management in `App.tsx`:
```tsx
const [isDrawerExpanded, setIsDrawerExpanded] = useState(true)
```

## Browser Support

PatternFly 5 supports:
- ✅ Chrome/Edge (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ⚠️ IE 11 (not supported)

## Development Tips

### Hot Reload
Changes to `.tsx` files reload automatically in dev mode.

### Component Explorer
Browse all PatternFly components at:
https://www.patternfly.org/components/all-components

### Icon Finder
Search 700+ icons at:
https://www.patternfly.org/design-foundations/icons

### Color Palette
See all colors at:
https://www.patternfly.org/design-foundations/colors

## Testing the UI

### Manual Test Checklist

- [ ] Header displays with red gradient
- [ ] Logo and title visible
- [ ] Drawer opens and closes
- [ ] Form fields accept input
- [ ] Helper text displays
- [ ] Execute button disabled when form invalid
- [ ] Execute button enabled when form valid
- [ ] Loading spinner shows during trace
- [ ] Error alert displays on failure
- [ ] Success status badge shows on success
- [ ] Graph renders in card
- [ ] Timeline shows flow steps
- [ ] Responsive layout works (resize window)
- [ ] Keyboard navigation works (tab through form)
- [ ] Focus indicators visible

### Sample Test Flow

1. **Open app** → See empty state
2. **Fill form**:
   - Source: `pod/web-pod`
   - Destination: `pod/db-pod`
   - Protocol: `tcp`
   - Dest Port: `5432`
3. **Click Execute Trace**
4. **See loading** → Spinner appears
5. **See results**:
   - Status badge (SUCCESS/DROPPED/ERROR)
   - Graph with nodes and edges
   - Timeline with steps
6. **Collapse drawer** → Click [x] button
7. **Expand drawer** → Click "Show Configuration"

## Next Steps

1. ✅ **Dependencies installed** → PatternFly packages
2. ✅ **UI redesigned** → All components use PatternFly
3. ✅ **Accessible** → WCAG 2.1 AA compliant
4. 🚧 **Test thoroughly** → Manual testing
5. 🚧 **Customize** → Adjust colors, spacing as needed
6. 🚧 **Add tests** → Component tests (future)

## Resources

- [PatternFly Get Started](https://www.patternfly.org/get-started/develop)
- [PatternFly React Docs](https://www.patternfly.org/components/all-components)
- [Red Hat UX](https://ux.redhat.com/)
- [PATTERNFLY_DESIGN.md](./PATTERNFLY_DESIGN.md) - Design system docs
- [PATTERNFLY_REDESIGN_SUMMARY.md](./PATTERNFLY_REDESIGN_SUMMARY.md) - What changed

---

**Ready to run**: `./run.sh dev` and visit http://localhost:3000
