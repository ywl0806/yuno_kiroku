import {useMutation} from '@tanstack/react-query';
import React, {useState} from 'react';
import {
  ActivityIndicator,
  Alert,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import colors from '../../colors';
import {useAuth} from '../context/AuthContext';
import {MyAxios} from '../lib/myAxios';

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
    backgroundColor: colors.main_bg,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 30,
    color: colors.main_blue,
  },
  input: {
    width: '100%',
    height: 50,
    borderWidth: 1,
    borderColor: colors.dark_gray,
    borderRadius: 8,
    paddingHorizontal: 15,
    marginBottom: 15,
    backgroundColor: 'white',
    fontSize: 16,
  },
  button: {
    width: '100%',
    height: 50,
    backgroundColor: colors.main_blue,
    borderRadius: 8,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 10,
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: 'bold',
  },
  buttonDisabled: {
    backgroundColor: colors.dark_gray,
  },
});

interface LoginResponse {
  message: string;
  user: {
    id: number;
    username: string;
    group_id: number;
    clan_group_id: number | null;
  };
  token: string;
}

export const LoginScreen = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const {login} = useAuth();

  const {mutate: handleLogin, isPending} = useMutation({
    mutationFn: async (data: {username: string; password: string}) => {
      const response = await MyAxios.post<LoginResponse>(
        '/api/auth/login',
        data,
      );
      return response.data;
    },
    onSuccess: async data => {
      try {
        await login(data.token);
      } catch (error) {
        Alert.alert('Error', 'Failed to save login information');
      }
    },
    onError: (error: any) => {
      const errorMessage =
        error.response?.data?.error || 'Login failed. Please try again.';
      Alert.alert('Login Error', errorMessage);
    },
  });

  const onSubmit = () => {
    if (!username.trim() || !password.trim()) {
      Alert.alert(
        'Validation Error',
        'Please enter both username and password',
      );
      return;
    }
    handleLogin({username: username.trim(), password});
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Login</Text>
      <TextInput
        style={styles.input}
        placeholder="Username"
        value={username}
        onChangeText={setUsername}
        autoCapitalize="none"
        editable={!isPending}
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        editable={!isPending}
        onSubmitEditing={onSubmit}
      />
      <TouchableOpacity
        style={[styles.button, isPending && styles.buttonDisabled]}
        onPress={onSubmit}
        disabled={isPending}>
        {isPending ? (
          <ActivityIndicator color="white" />
        ) : (
          <Text style={styles.buttonText}>Login</Text>
        )}
      </TouchableOpacity>
    </View>
  );
};
