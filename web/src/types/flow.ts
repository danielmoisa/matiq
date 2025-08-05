export enum NodeType {
  // Virtual/Local Actions
  TRANSFORMER = 'transformer',
  
  // APIs
  RESTAPI = 'restapi',
  GRAPHQL = 'graphql',
  
  // Cache/Messaging
  REDIS = 'redis',
  UPSTASH = 'upstash',
  
  // Databases
  MYSQL = 'mysql',
  MARIADB = 'mariadb',
  POSTGRESQL = 'postgresql',
  MONGODB = 'mongodb',
  TIDB = 'tidb',
  ELASTICSEARCH = 'elasticsearch',
  SUPABASEDB = 'supabasedb',
  FIREBASE = 'firebase',
  CLICKHOUSE = 'clickhouse',
  MSSQL = 'mssql',
  DYNAMODB = 'dynamodb',
  SNOWFLAKE = 'snowflake',
  COUCHDB = 'couchdb',
  ORACLE = 'oracle',
  ORACLE_9I = 'oracle9i',
  NEON = 'neon',
  HYDRA = 'hydra',
  
  // Storage
  S3 = 's3',
  
  // Communication
  SMTP = 'smtp',
  
  // AI/ML
  HUGGINGFACE = 'huggingface',
  HFENDPOINT = 'hfendpoint',
  AI_AGENT = 'aiagent',
  
  // External Services
  GOOGLESHEETS = 'googlesheets',
  AIRTABLE = 'airtable',
  APPWRITE = 'appwrite',
  
  // Flow Control
  TRIGGER = 'trigger',
  SERVER_SIDE_TRANSFORMER = 'serversidetransformer',
  CONDITION = 'condition',
  WEBHOOK_RESPONSE = 'webhookresponse',
  WF_DRIVE = 'wfdrive',
  
  // Legacy - keeping for backward compatibility
  SCHEDULE = 'schedule',
  WEBHOOK = 'webhook',
  REST_API = 'rest-api',
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

export interface FlowNode {
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

// Backend node structure from API response
interface BackendNode {
  id?: string;
  'action_type,'?: string;
  action_type?: string;
  type?: string;
  triggerType?: string;
  position?: { x: number; y: number };
  data?: Record<string, unknown>;
  connections?: string[];
}

// Backend response structure (matches actual Go backend response)
export interface FlowBackend {
  uid: string;                  // UUID from backend
  resourceID?: string;          // Encoded resource ID (optional field)
  displayName: string;          // Flow name
  actionType: string;           // Flow type as string
  template: Record<string, unknown>;  // Raw template data from backend (map[string]interface{})
  transformer: unknown;         // Transformer data (can be null)
  triggerMode: string;          // Trigger mode as string
  config: unknown;              // Config data (can be null)
  createdAt: string;            // ISO timestamp
  createdBy: string;            // User ID who created
  updatedAt: string;            // ISO timestamp
  updatedBy: string;            // User ID who updated
}

// Frontend-friendly flow interface (converted from backend)
export interface Flow {
  id: string;                   // Use resourceID as string
  name: string;
  description?: string;
  nodes: FlowNode[];
  connections: Connection[];
  isActive: boolean;
  status: "active" | "draft" | "paused" | "error";
  triggerMode?: string;
  createdAt: string;
  updatedAt: string;
  // Backend fields for API calls
  uid?: string;
  resourceID?: string;
  actionType?: string;
}

// Template structure stored in backend template field
interface FlowTemplate {
  nodes?: FlowNode[];
  connections?: Connection[];
  resourceID?: number;
  runByAnonymous?: boolean;
  teamID?: number;
}


// Convert backend flow to frontend flow
export function convertBackendToFrontend(backend: FlowBackend | null | undefined): Flow {
  // Handle undefined/null backend response
  if (!backend) {
    return {
      id: 'unknown',
      name: 'Unknown Flow',
      description: '',
      nodes: [],
      connections: [],
      isActive: false,
      status: 'error',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
  }

  // Extract nodes and connections from template
  let nodes: FlowNode[] = [];
  let connections: Connection[] = [];

  if (backend.template && typeof backend.template === 'object') {
    // Safely access template data as it comes from backend as map[string]interface{}
    const templateData = backend.template;
    
    // Transform backend nodes to frontend format
    if (Array.isArray(templateData.nodes)) {
      nodes = templateData.nodes.map((backendNode: BackendNode) => {
        // Handle the backend node structure which uses action_type, (with comma)
        const nodeType = backendNode['action_type,'] || backendNode.action_type || backendNode.type || 'unknown';
        
        return {
          id: backendNode.id || 'unknown',
          type: nodeType,
          triggerType: backendNode.triggerType,
          position: backendNode.position || { x: 0, y: 0 },
          data: backendNode.data || {},
          connections: backendNode.connections || []
        } as FlowNode;
      });
    }
    
    connections = Array.isArray(templateData.connections) ? templateData.connections as Connection[] : [];
  }

  // Determine status based on template and trigger mode
  let status: "active" | "draft" | "paused" | "error" = "draft";
  if (nodes.length > 0) {
    status = "active";
  }

  return {
    id: backend.resourceID || backend.uid || 'unknown', 
    name: backend.displayName || 'Untitled Flow',
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
    resourceID: backend.resourceID,
    actionType: backend.actionType,
  };
}

// Convert frontend flow to backend create/update request
export function convertFrontendToBackendRequest(
  flow: Partial<Flow>,
  nodes: FlowNode[] = [],
  connections: Connection[] = []
): {
  displayName: string;
  actionType: string;
  triggerMode: string;
  template: FlowTemplate;
  transformer?: unknown;
  config?: unknown;
} {
  // Create template object from nodes and connections
  const template: FlowTemplate = {
    nodes,
    connections,
    runByAnonymous: true,
  };

  return {
    displayName: flow.name || "Untitled Flow",
    actionType: flow.actionType || "restapi",
    triggerMode: flow.triggerMode || "1", // Backend expects string
    template,
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
