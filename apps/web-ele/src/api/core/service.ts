import { requestClient } from '#/api/request';

export interface OpenListService {
  id?: number;
  serviceName: string;
  account: string;
  token: string;
  serviceUrl: string;
  backupUrl?: string;
  enabled: boolean;
}

export function getServiceListApi(): Promise<OpenListService[]> {
  return requestClient.get<OpenListService[]>('/openlist/service');
}

export function createServiceApi(data: OpenListService): Promise<void> {
  return requestClient.post('/openlist/service', data);
}

export function updateServiceApi(
  id: number,
  data: OpenListService,
): Promise<void> {
  return requestClient.put(`/openlist/service/${id}`, data);
}

export function deleteServiceApi(id: number): Promise<void> {
  return requestClient.delete(`/openlist/service/${id}`);
}
