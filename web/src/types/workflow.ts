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

export interface Workflow {
  id: string;
  name: string;
  description?: string;
  nodes: WorkflowNode[];
  connections: Connection[];
  isActive: boolean;
  status: 'active' | 'draft' | 'paused' | 'error';
  createdAt: string;
  updatedAt: string;
}
