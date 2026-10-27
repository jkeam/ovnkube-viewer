import * as React from 'react';
import { Page, PageSection, PageSectionVariants } from '@patternfly/react-core';
import App from '../App';

/**
 * OVN Trace Viewer Console Plugin Page
 * This component wraps the main App to integrate with OpenShift Console
 */
const OVNTraceViewerPage: React.FC = () => {
  return (
    <Page>
      <PageSection variant={PageSectionVariants.light}>
        <App />
      </PageSection>
    </Page>
  );
};

export default OVNTraceViewerPage;
