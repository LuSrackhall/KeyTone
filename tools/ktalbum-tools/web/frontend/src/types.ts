export interface FileInfo {
  name: string;
  version: number;
  exportTime: string;
  albumUUID: string;
}

export interface ApiResponse<T = any> {
  message: string;
  data?: T;
  error?: string;
}

export interface UploadResponse extends ApiResponse {
  message: "ok" | string;
}

export interface FileInfoResponse extends ApiResponse {
  message: "ok" | string;
  data?: FileInfo;
}
