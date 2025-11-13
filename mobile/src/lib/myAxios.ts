import AsyncStorage from '@react-native-async-storage/async-storage';
import axios from 'axios';

import {logoutFromStorage} from '../context/AuthContext';

export const MyAxios = axios.create({
  baseURL: 'http://localhost:1323',
});

const TOKEN_KEY = '@auth_token';

// Request interceptor to add token to headers
MyAxios.interceptors.request.use(
  async config => {
    try {
      const token = await AsyncStorage.getItem(TOKEN_KEY);
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    } catch (error) {
      console.error('Error getting token:', error);
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

// Response interceptor to handle 401 errors
MyAxios.interceptors.response.use(
  response => {
    return response;
  },
  async error => {
    if (error.response?.status === 401) {
      // Token expired or invalid, logout user
      await logoutFromStorage();
    }
    return Promise.reject(error);
  },
);
