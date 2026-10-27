import { useState, FormEvent } from 'react'
import {
  Form,
  FormGroup,
  FormSection,
  TextInput,
  FormSelect,
  FormSelectOption,
  Button,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core'
import { PlayIcon } from '@patternfly/react-icons'

interface TraceFormProps {
  onSubmit: (request: {
    source: string
    destination: string
    protocol: string
    sourcePort?: number
    destPort?: number
  }) => void
  loading: boolean
}

function TraceForm({ onSubmit, loading }: TraceFormProps) {
  const [source, setSource] = useState('')
  const [destination, setDestination] = useState('')
  const [protocol, setProtocol] = useState('tcp')
  const [sourcePort, setSourcePort] = useState('')
  const [destPort, setDestPort] = useState('')

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()

    const request = {
      source,
      destination,
      protocol,
      ...(sourcePort && { sourcePort: parseInt(sourcePort, 10) }),
      ...(destPort && { destPort: parseInt(destPort, 10) }),
    }

    onSubmit(request)
  }

  const protocolOptions = [
    { value: 'tcp', label: 'TCP', disabled: false },
    { value: 'udp', label: 'UDP', disabled: false },
    { value: 'icmp', label: 'ICMP', disabled: false },
  ]

  return (
    <Form onSubmit={handleSubmit}>
      <FormSection title="Network Endpoints">
        <FormGroup label="Source" isRequired fieldId="source">
          <TextInput
            isRequired
            type="text"
            id="source"
            name="source"
            value={source}
            onChange={(_event, value) => setSource(value)}
            placeholder="e.g., pod/web-pod or 10.128.0.5"
            isDisabled={loading}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>Pod name, IP address, or namespace/pod</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>

        <FormGroup label="Destination" isRequired fieldId="destination">
          <TextInput
            isRequired
            type="text"
            id="destination"
            name="destination"
            value={destination}
            onChange={(_event, value) => setDestination(value)}
            placeholder="e.g., svc/backend or 10.96.0.1"
            isDisabled={loading}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>Service, Pod, IP, or external address</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>
      </FormSection>

      <FormSection title="Protocol & Ports">
        <FormGroup label="Protocol" isRequired fieldId="protocol">
          <FormSelect
            value={protocol}
            onChange={(_event, value) => setProtocol(value as string)}
            id="protocol"
            name="protocol"
            isDisabled={loading}
          >
            {protocolOptions.map((option, index) => (
              <FormSelectOption
                key={index}
                value={option.value}
                label={option.label}
                isDisabled={option.disabled}
              />
            ))}
          </FormSelect>
        </FormGroup>

        <FormGroup label="Source Port" fieldId="sourcePort">
          <TextInput
            type="number"
            id="sourcePort"
            name="sourcePort"
            value={sourcePort}
            onChange={(_event, value) => setSourcePort(value)}
            placeholder="Optional"
            min={1}
            max={65535}
            isDisabled={loading}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>Optional source port (1-65535)</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>

        <FormGroup label="Destination Port" fieldId="destPort">
          <TextInput
            type="number"
            id="destPort"
            name="destPort"
            value={destPort}
            onChange={(_event, value) => setDestPort(value)}
            placeholder="e.g., 8080"
            min={1}
            max={65535}
            isDisabled={loading}
          />
          <FormHelperText>
            <HelperText>
              <HelperTextItem>Destination port (1-65535)</HelperTextItem>
            </HelperText>
          </FormHelperText>
        </FormGroup>
      </FormSection>

      <Button
        type="submit"
        variant="primary"
        isBlock
        isDisabled={loading || !source || !destination}
        icon={<PlayIcon />}
        isLoading={loading}
        spinnerAriaValueText={loading ? 'Executing trace...' : undefined}
      >
        {loading ? 'Executing Trace...' : 'Execute Trace'}
      </Button>
    </Form>
  )
}

export default TraceForm
