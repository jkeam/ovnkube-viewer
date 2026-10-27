import { useState } from 'react'
import {
  Page,
  PageSection,
  PageSectionVariants,
  Masthead,
  MastheadMain,
  MastheadBrand,
  MastheadContent,
  Brand,
  Drawer,
  DrawerContent,
  DrawerContentBody,
  DrawerPanelContent,
  DrawerHead,
  DrawerActions,
  DrawerCloseButton,
  Title,
  EmptyState,
  EmptyStateIcon,
  EmptyStateBody,
  EmptyStateHeader,
  Spinner,
  Alert,
  AlertVariant,
  AlertActionCloseButton,
} from '@patternfly/react-core'
import { TopologyIcon } from '@patternfly/react-icons'
import TraceForm from './components/TraceForm/TraceForm'
import FlowVisualization from './components/FlowVisualization/FlowVisualization'
import { TraceResult } from './types'

function App() {
  const [traceResult, setTraceResult] = useState<TraceResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [isDrawerExpanded, setIsDrawerExpanded] = useState(true)

  const handleTraceSubmit = async (request: {
    source: string
    destination: string
    protocol: string
    sourcePort?: number
    destPort?: number
  }) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/v1/trace', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data: TraceResult = await response.json()
      setTraceResult(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
      console.error('Error executing trace:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCloseError = () => {
    setError(null)
  }

  const panelContent = (
    <DrawerPanelContent minSize="320px" defaultSize="350px" maxSize="450px">
      <DrawerHead>
        <Title headingLevel="h2" size="xl">
          Trace Configuration
        </Title>
        <DrawerActions>
          <DrawerCloseButton onClick={() => setIsDrawerExpanded(false)} />
        </DrawerActions>
      </DrawerHead>
      <DrawerContentBody>
        <TraceForm onSubmit={handleTraceSubmit} loading={loading} />
      </DrawerContentBody>
    </DrawerPanelContent>
  )

  const pageHeader = (
    <Masthead style={{ background: 'linear-gradient(135deg, #ee0000 0%, #c00 100%)' }}>
      <MastheadMain>
        <MastheadBrand>
          <Brand
            src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23ffffff' d='M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z'/%3E%3C/svg%3E"
            alt="OVN Logo"
            widths={{ default: '40px' }}
          />
        </MastheadBrand>
        <div style={{ marginLeft: '1rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <Title headingLevel="h1" size="2xl" style={{ color: '#fff', margin: 0 }}>
            OVN-Kubernetes Trace Viewer
          </Title>
          <span style={{ color: 'rgba(255, 255, 255, 0.9)', fontSize: '0.875rem', marginTop: '0.25rem' }}>
            Visualize network traffic flow through OVN components
          </span>
        </div>
      </MastheadMain>
    </Masthead>
  )

  return (
    <Page header={pageHeader} isManagedSidebar>
      <Drawer isExpanded={isDrawerExpanded} isInline>
        <DrawerContent panelContent={panelContent}>
          <DrawerContentBody className="main-drawer-content">
            {error && (
              <Alert
                variant={AlertVariant.danger}
                title="Trace Execution Error"
                actionClose={<AlertActionCloseButton onClose={handleCloseError} />}
                style={{ margin: '1rem' }}
              >
                {error}
              </Alert>
            )}

            {loading && (
              <div style={{ height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <EmptyState>
                  <EmptyStateHeader
                    titleText="Executing trace..."
                    icon={<EmptyStateIcon icon={Spinner} />}
                    headingLevel="h4"
                  />
                  <EmptyStateBody>
                    Please wait while we trace the network path through OVN components.
                  </EmptyStateBody>
                </EmptyState>
              </div>
            )}

            {traceResult && !loading && (
              <FlowVisualization
                result={traceResult}
                onToggleDrawer={() => setIsDrawerExpanded(!isDrawerExpanded)}
                isDrawerOpen={isDrawerExpanded}
              />
            )}

            {!traceResult && !loading && !error && (
              <div style={{ height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <EmptyState>
                  <EmptyStateHeader
                    titleText="Ready to trace"
                    icon={<EmptyStateIcon icon={TopologyIcon} />}
                    headingLevel="h4"
                  />
                  <EmptyStateBody>
                    Configure trace parameters in the side panel and click "Execute Trace" to
                    visualize the network flow through OVN-Kubernetes components.
                    <br />
                    <br />
                    The visualization will show the packet path through logical switches, routers,
                    ACLs, and other OVN elements.
                  </EmptyStateBody>
                </EmptyState>
              </div>
            )}
          </DrawerContentBody>
        </DrawerContent>
      </Drawer>
    </Page>
  )
}

export default App
