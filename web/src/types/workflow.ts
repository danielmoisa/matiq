export enum NodeType {
  // Triggers
  SCHEDULE = 'schedule',
  WEBHOOK = 'webhook',
  
  // Databases
  POSTGRES = 'postgres',
  MYSQL = 'mysql',
  MARIADB = 'mariadb',
  TIDB = 'tidb',
  NEON = 'neon',
  MONGODB = 'mongodb',
  SNOWFLAKE = 'snowflake',
  SUPABASE = 'supabase',
  CLICKHOUSE = 'clickhouse',
  HYDRA = 'hydra',
  
  // APIs
  REST_API = 'rest-api',
  GRAPHQL = 'graphql',
  
  // Actions
  AI_AGENT = 'ai-agent',
  TRANSFORMER = 'transformer',
  CONDITION = 'condition',
  LOOP = 'loop',
  RESPONSE = 'response',
  ERROR_HANDLER = 'error-handler'
}

export type TriggerType = 'schedule' | 'webhook';

export type EventType = NodeType;

export interface Position {
  x: number;
  y: number;
}

export interface WorkflowNode {
  id: string;
  type: EventType;
  triggerType?: TriggerType;
  position: Position;
  data: Record<string, unknown>;
  connections?: string[];
}

export interface Connection {
  id: string;
  sourceId: string;
  targetId: string;
}

// Backend response structure (matches Go model)
export interface WorkflowBackend {
  id: number;                    // Backend primary key
  uid: string;                   // UUID from backend
  teamId: number;               // Team ID
  workflowId: number;           // Workflow ID within team
  version: number;              // Version control
  resourceId: number;           // Resource ID
  name: string;                 // Workflow name (DisplayName in backend)
  type: number;                 // Workflow type as integer
  triggerMode: string;          // "schedule" | "webhook"
  transformer: string;          // JSON string
  template: string;             // JSON string containing nodes/connections
  config: string;               // JSON string
  createdAt: string;            // ISO timestamp
  createdBy: number;            // User ID
  updatedAt: string;            // ISO timestamp
  updatedBy: number;            // User ID
}

// Frontend-friendly workflow interface (converted from backend)
export interface Workflow {
  id: string;                   // Use workflowId as string
  name: string;
  description?: string;
  nodes: WorkflowNode[];
  connections: Connection[];
  isActive: boolean;
  status: 'active' | 'draft' | 'paused' | 'error';
  triggerMode?: string;
  createdAt: string;
  updatedAt: string;
  // Backend fields for API calls
  teamId?: number;
  workflowId?: number;
  version?: number;
  resourceId?: number;
  type?: number;
}

// Helper functions to convert between backend and frontend models

// Template structure stored in backend's template JSON field
interface WorkflowTemplate {
  nodes?: WorkflowNode[];
  connections?: Connection[];
  resourceID?: number;
  runByAnonymous?: boolean;
  teamID?: number;
}

// Convert backend workflow to frontend workflow
export function convertBackendToFrontend(backend: WorkflowBackend): Workflow {
  // Parse template JSON to extract nodes and connections
  let nodes: WorkflowNode[] = [];
  let connections: Connection[] = [];
  
  if (backend.template) {
    try {
      const template: WorkflowTemplate = JSON.parse(backend.template);
      nodes = template.nodes || [];
      connections = template.connections || [];
    } catch (error) {
      console.warn('Failed to parse workflow template:', error);
    }
  }

  // Determine status based on template content and trigger mode
  let status: 'active' | 'draft' | 'paused' | 'error' = 'draft';
  if (nodes.length > 0) {
    status = 'active';
  }

  return {
    id: backend.workflowId.toString(),
    name: backend.name,
    description: '', // Backend doesn't have description field
    nodes,
    connections,
    isActive: nodes.length > 0, // Active if has nodes
    status,
    triggerMode: backend.triggerMode,
    createdAt: backend.createdAt,
    updatedAt: backend.updatedAt,
    // Keep backend fields for API calls
    teamId: backend.teamId,
    workflowId: backend.workflowId,
    version: backend.version,
    resourceId: backend.resourceId,
    type: backend.type,
  };
}

// Convert frontend workflow to backend create/update request
export function convertFrontendToBackendRequest(
  workflow: Partial<Workflow>,
  nodes: WorkflowNode[] = [],
  connections: Connection[] = []
): {
  displayName: string;
  triggerMode: string;
  template: string;
  transformer: string;
  config: string;
  resourceId?: string;
  workflowType?: string;
} {
  // Create template JSON from nodes and connections
  const template: WorkflowTemplate = {
    nodes,
    connections,
  };

  return {
    displayName: workflow.name || 'Untitled Workflow',
    triggerMode: workflow.triggerMode || 'webhook',
    template: JSON.stringify(template),
    transformer: '{}', // Empty transformer for now
    config: '{}', // Empty config for now
    resourceId: workflow.resourceId?.toString() || '0',
    workflowType: 'workflow', // Default type
  };
}

// API Response wrapper (matches backend response format)
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface ApiError {
  success: false;
  error: string;
  message?: string;
}
