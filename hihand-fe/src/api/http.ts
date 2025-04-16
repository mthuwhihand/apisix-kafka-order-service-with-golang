// src/api/http.ts
import axios from 'axios';

const http = axios.create({
    baseURL: 'http://localhost:9080',
    timeout: 5000,
});

export default http;
