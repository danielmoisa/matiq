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

// Backend response structure (matches actual Go backend response)
export interface WorkflowBackend {
  workflowID: string;           // Main workflow ID (from JSON tag)
  uid: string;                  // UUID from backend
  teamID: string;               // Encoded team ID
  version: number;              // Version number
  resourceID?: string;          // Encoded resource ID (optional field)
  displayName: string;          // Workflow name
  workflowType: string;         // Workflow type as string
  isVirtualResource: boolean;   // Virtual resource flag
  content: {                    // Content containing nodes and connections
    nodes?: WorkflowNode[];
    connections?: Connection[];
    resourceID?: number;
    runByAnonymous?: boolean;
    teamID?: number;
  };
  transformer: unknown;         // Transformer data (can be null)
  triggerMode: string;          // Trigger mode as string
  template: unknown;            // Template data (can be null)
  config: unknown;              // Config data (can be null)
  createdAt: string;            // ISO timestamp
  createdBy: string;            // Encoded user ID
  updatedAt: string;            // ISO timestamp
  updatedBy: string;            // Encoded user ID
}

// Frontend-friendly workflow interface (converted from backend)
export interface Workflow {
  id: string;                   // Use resourceID as string
  name: string;
  description?: string;
  nodes: WorkflowNode[];
  connections: Connection[];
  isActive: boolean;
  status: "active" | "draft" | "paused" | "error";
  triggerMode?: string;
  createdAt: string;
  updatedAt: string;
  // Backend fields for API calls
  uid?: string;
  teamID?: string;
  version?: number;
  resourceID?: string;
  workflowType?: string;
}

// Template structure stored in backend content field
interface WorkflowContent {
  nodes?: WorkflowNode[];
  connections?: Connection[];
  resourceID?: number;
  runByAnonymous?: boolean;
  teamID?: number;
}

// Convert backend workflow to frontend workflow
export function convertBackendToFrontend(backend: WorkflowBackend | null | undefined): Workflow {
  // Handle undefined/null backend response
  if (!backend) {
    return {
      id: 'unknown',
      name: 'Unknown Workflow',
      description: '',
      nodes: [],
      connections: [],
      isActive: false,
      status: 'error',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
  }

  // Extract nodes and connections from content
  let nodes: WorkflowNode[] = [];
  let connections: Connection[] = [];
  
  if (backend.content) {
    nodes = backend.content.nodes || [];
    connections = backend.content.connections || [];
  }

  // Determine status based on content and trigger mode
  let status: "active" | "draft" | "paused" | "error" = "draft";
  if (nodes.length > 0) {
    status = "active";
  }

  return {
    id: backend.workflowID || backend.resourceID || backend.uid || 'unknown', // Use workflowID first
    name: backend.displayName || 'Untitled Workflow',
    description: "", // Backend does not have description field
    nodes,
    connections,
    isActive: nodes.length > 0, // Active if has nodes
    status,
    triggerMode: backend.triggerMode,
    createdAt: backend.createdAt || new Date().toISOString(),
    updatedAt: backend.updatedAt || new Date().toISOString(),
    // Keep backend fields for API calls
    uid: backend.uid,
    teamID: backend.teamID,
    version: backend.version,
    resourceID: backend.resourceID,
    workflowType: backend.workflowType,
  };
}

// Convert frontend workflow to backend create/update request
export function convertFrontendToBackendRequest(
  workflow: Partial<Workflow>,
  nodes: WorkflowNode[] = [],
  connections: Connection[] = []
): {
  displayName: string;
  workflowType: string;
  triggerMode: string;
  content: WorkflowContent;
  transformer?: unknown;
  config?: unknown;
} {
  // Create content object from nodes and connections
  const content: WorkflowContent = {
    nodes,
    connections,
    runByAnonymous: true,
  };

  return {
    displayName: workflow.name || "Untitled Workflow",
    workflowType: workflow.workflowType || "restapi",
    triggerMode: workflow.triggerMode || "1", // Backend expects string
    content,
    transformer: null,
    config: null,
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
