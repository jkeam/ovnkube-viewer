import {
  Card,
  CardBody,
  CardTitle,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Button,
  Label,
  Split,
  SplitItem,
  Badge,
  Flex,
  FlexItem,
  Text,
  TextContent,
  TextVariants,
  Stack,
  StackItem,
  Grid,
  GridItem,
  List,
  ListItem,
} from '@patternfly/react-core'
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
  TimesCircleIcon,
  ArrowRightIcon,
  BarsIcon,
  NetworkIcon,
} from '@patternfly/react-icons'
import { TraceResult } from '../../types'
import './FlowVisualization.css'

interface FlowVisualizationProps {
  result: TraceResult
  onToggleDrawer: () => void
  isDrawerOpen: boolean
}

function FlowVisualization({ result, onToggleDrawer, isDrawerOpen }: FlowVisualizationProps) {
  const getStatusIcon = () => {
    switch (result.metadata.status) {
      case 'success':
        return <CheckCircleIcon color="var(--pf-v5-global--success-color--100)" />
      case 'dropped':
        return <TimesCircleIcon color="var(--pf-v5-global--danger-color--100)" />
      case 'error':
        return <ExclamationCircleIcon color="var(--pf-v5-global--warning-color--100)" />
      default:
        return null
    }
  }

  const getStatusColor = () => {
    switch (result.metadata.status) {
      case 'success':
        return 'green'
      case 'dropped':
        return 'red'
      case 'error':
        return 'orange'
      default:
        return 'grey'
    }
  }

  const getNodeColor = (type: string): string => {
    switch (type) {
      case 'switch':
        return 'var(--pf-v5-global--primary-color--100)'
      case 'router':
        return 'var(--pf-v5-global--danger-color--100)'
      case 'port':
        return 'var(--pf-v5-global--palette--purple-400)'
      case 'acl':
        return 'var(--pf-v5-global--warning-color--100)'
      default:
        return 'var(--pf-v5-global--Color--200)'
    }
  }

  return (
    <div className="flow-visualization-container">
      <Toolbar isSticky>
        <ToolbarContent>
          {!isDrawerOpen && (
            <ToolbarItem>
              <Button
                variant="plain"
                onClick={onToggleDrawer}
                icon={<BarsIcon />}
                aria-label="Toggle drawer"
              >
                Show Configuration
              </Button>
            </ToolbarItem>
          )}
          <ToolbarItem>
            <Split hasGutter>
              <SplitItem>
                <Label icon={getStatusIcon()} color={getStatusColor() as any}>
                  {result.metadata.status.toUpperCase()}
                </Label>
              </SplitItem>
              <SplitItem>
                <Text component={TextVariants.small}>
                  <strong>Source:</strong> {result.metadata.source}
                </Text>
              </SplitItem>
              <SplitItem>
                <ArrowRightIcon />
              </SplitItem>
              <SplitItem>
                <Text component={TextVariants.small}>
                  <strong>Destination:</strong> {result.metadata.destination}
                </Text>
              </SplitItem>
              <SplitItem>
                <Badge>{result.metadata.protocol.toUpperCase()}</Badge>
              </SplitItem>
              <SplitItem>
                <Badge isRead>{result.flow.length} steps</Badge>
              </SplitItem>
            </Split>
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>

      <div className="flow-content-area">
        <Grid hasGutter>
          <GridItem span={12} lg={6}>
            <Card isFullHeight>
              <CardTitle>
                <Flex spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>
                    <NetworkIcon /> Network Topology
                  </FlexItem>
                  <FlexItem>
                    <Badge isRead>{result.topology.nodes.length} nodes</Badge>
                  </FlexItem>
                </Flex>
              </CardTitle>
              <CardBody>
                <Stack hasGutter>
                  {result.topology.nodes.map((node, index) => (
                    <StackItem key={node.id}>
                      <Card isCompact className="topology-node-card">
                        <CardBody>
                          <Flex
                            alignItems={{ default: 'alignItemsCenter' }}
                            spaceItems={{ default: 'spaceItemsSm' }}
                          >
                            <FlexItem>
                              <div
                                className="node-indicator"
                                style={{ backgroundColor: getNodeColor(node.type) }}
                              />
                            </FlexItem>
                            <FlexItem flex={{ default: 'flex_1' }}>
                              <TextContent>
                                <Text component={TextVariants.h4}>{node.name}</Text>
                                <Text component={TextVariants.small}>
                                  <Badge>{node.type}</Badge>
                                </Text>
                              </TextContent>
                            </FlexItem>
                            {index < result.topology.nodes.length - 1 && (
                              <FlexItem>
                                <ArrowRightIcon color="var(--pf-v5-global--Color--200)" />
                              </FlexItem>
                            )}
                          </Flex>
                        </CardBody>
                      </Card>
                    </StackItem>
                  ))}
                </Stack>

                {result.topology.edges.length > 0 && (
                  <>
                    <br />
                    <TextContent>
                      <Text component={TextVariants.h4}>Connections</Text>
                    </TextContent>
                    <List isPlain>
                      {result.topology.edges.map((edge, idx) => (
                        <ListItem key={idx}>
                          {edge.source} <ArrowRightIcon /> {edge.target}
                          {edge.label && ` (${edge.label})`}
                        </ListItem>
                      ))}
                    </List>
                  </>
                )}
              </CardBody>
            </Card>
          </GridItem>

          <GridItem span={12} lg={6}>
            <Card isFullHeight>
              <CardTitle>
                <Flex spaceItems={{ default: 'spaceItemsSm' }}>
                  <FlexItem>Flow Steps</FlexItem>
                  <FlexItem>
                    <Badge>{result.flow.length}</Badge>
                  </FlexItem>
                </Flex>
              </CardTitle>
              <CardBody>
                <Stack hasGutter>
                  {result.flow.map((step) => (
                    <StackItem key={step.stepNum}>
                      <Card isCompact className="flow-step-card">
                        <CardBody>
                          <Flex direction={{ default: 'column' }}>
                            <FlexItem>
                              <Flex
                                alignItems={{ default: 'alignItemsCenter' }}
                                spaceItems={{ default: 'spaceItemsSm' }}
                              >
                                <FlexItem>
                                  <Badge className="step-number">{step.stepNum}</Badge>
                                </FlexItem>
                                <FlexItem>
                                  <Text component={TextVariants.h4}>
                                    <Badge color="blue">{step.action}</Badge>
                                  </Text>
                                </FlexItem>
                              </Flex>
                            </FlexItem>
                            <FlexItem>
                              <Text component={TextVariants.p}>{step.description}</Text>
                            </FlexItem>
                            <FlexItem>
                              <Text
                                component={TextVariants.small}
                                className="packet-info"
                              >
                                {step.packet.sourceIP} → {step.packet.destIP} (
                                {step.packet.protocol}
                                {step.packet.sourcePort && `:${step.packet.sourcePort}`}
                                {step.packet.destPort && ` → :${step.packet.destPort}`})
                              </Text>
                            </FlexItem>
                          </Flex>
                        </CardBody>
                      </Card>
                    </StackItem>
                  ))}
                </Stack>
              </CardBody>
            </Card>
          </GridItem>
        </Grid>
      </div>
    </div>
  )
}

export default FlowVisualization
