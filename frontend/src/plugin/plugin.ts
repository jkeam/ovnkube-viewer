import {
  Plugin,
  ModelDefinition,
  ModelFeatureFlag,
  RoutePage,
  HrefNavItem,
} from '@openshift-console/dynamic-plugin-sdk';

type ConsumedExtensions =
  | ModelDefinition
  | ModelFeatureFlag
  | RoutePage
  | HrefNavItem;

const plugin: Plugin<ConsumedExtensions> = {
  name: 'ovnkube-viewer',
  version: '0.1.0',
  displayName: 'OVN Trace Viewer',
  description: 'Visual viewer for OVN-Kubernetes traffic tracing',
  exposedModules: {
    OVNTraceViewerPage: './OVNTraceViewerPage',
  },
  dependencies: {
    '@console/pluginAPI': '*',
  },
};

export default plugin;
