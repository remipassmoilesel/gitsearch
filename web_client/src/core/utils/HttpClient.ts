import axios, {AxiosInstance} from 'axios';

export type HttpClient = AxiosInstance;

export const HttpClient = axios.create({
    baseURL: '/api',
    timeout: 30_000, // High timeout for indexing operations
    headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
    },
});
