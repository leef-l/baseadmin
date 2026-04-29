/** 守护进程类型定义 */

export interface DaemonItem {
  id: string;
  name: string;
  program: string;
  command: string;
  directory: string;
  runUser: string;
  numprocs: number;
  priority: number;
  autostart: number;
  autorestart: number;
  startsecs: number;
  startretries: number;
  stopSignal: string;
  environment?: string;
  remark?: string;
  configPath?: string;
  outLogPath?: string;
  errLogPath?: string;
  runStatus?: string;
  pid?: string;
  uptime?: string;
  statusText?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface DaemonListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  program?: string;
}

export interface DaemonCreateParams {
  name: string;
  program: string;
  command: string;
  directory: string;
  runUser?: string;
  numprocs?: number;
  priority?: number;
  autostart?: number;
  autorestart?: number;
  startsecs?: number;
  startretries?: number;
  stopSignal?: string;
  environment?: string;
  remark?: string;
}

export interface DaemonUpdateParams extends DaemonCreateParams {
  id: string;
}

export interface DaemonOperationResult {
  program: string;
  runStatus: string;
  message: string;
}

export interface DaemonBatchOperationResult {
  results: DaemonOperationResult[];
}

export interface DaemonLogResult {
  program: string;
  logType: string;
  content: string;
}
