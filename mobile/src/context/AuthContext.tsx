import AsyncStorage from '@react-native-async-storage/async-storage';
import React, {createContext, useContext, useEffect, useState} from 'react';
import {AppState, AppStateStatus} from 'react-native';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
  login: (token: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const TOKEN_KEY = '@auth_token';

// Global logout function that can be called from anywhere
let globalLogout: (() => Promise<void>) | null = null;

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    checkAuthStatus();

    // Listen for app state changes to check token status
    const subscription = AppState.addEventListener(
      'change',
      (nextAppState: AppStateStatus) => {
        if (nextAppState === 'active') {
          checkAuthStatus();
        }
      },
    );

    return () => {
      subscription.remove();
    };
  }, []);

  const checkAuthStatus = async () => {
    try {
      const storedToken = await AsyncStorage.getItem(TOKEN_KEY);
      if (storedToken) {
        setToken(storedToken);
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error('Error checking auth status:', error);
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (newToken: string) => {
    try {
      await AsyncStorage.setItem(TOKEN_KEY, newToken);
      setToken(newToken);
      setIsAuthenticated(true);
    } catch (error) {
      console.error('Error saving token:', error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await AsyncStorage.removeItem(TOKEN_KEY);
      setToken(null);
      setIsAuthenticated(false);
    } catch (error) {
      console.error('Error removing token:', error);
      throw error;
    }
  };

  // Set global logout function
  useEffect(() => {
    globalLogout = logout;
    return () => {
      globalLogout = null;
    };
  }, []);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        token,
        login,
        logout,
      }}>
      {children}
    </AuthContext.Provider>
  );
};

// Export function to logout from anywhere (e.g., from axios interceptor)
export const logoutFromStorage = async () => {
  if (globalLogout) {
    await globalLogout();
  } else {
    // Fallback: directly remove token from storage
    try {
      await AsyncStorage.removeItem(TOKEN_KEY);
    } catch (error) {
      console.error('Error removing token:', error);
    }
  }
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

