import { env } from '@/utils/env';

export async function apiRequest<T = any>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${env.apiUrl}${endpoint}`, options);
  if (!response.ok) throw new Error('API Error');
  return response.json();
}
