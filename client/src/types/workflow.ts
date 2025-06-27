export type TriggerType = 'schedule' | 'webhook';

export type EventType = 
  | 'schedule'
  | 'webhook'
  | 'postgres' 
  | 'mysql' 
  | 'mariadb' 
  | 'tidb' 
  | 'neon' 
  | 'mongodb' 
  | 'snowflake' 
  | 'supabase' 
  | 'clickhouse' 
  | 'hydra'
  | 'rest-api'
  | 'graphql'
  | 'ai-agent'
  | 'transformer'
  | 'condition'
  | 'loop'
  | 'response'
  | 'error-handler';

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
  createdAt: Date;
  updatedAt: Date;
}
