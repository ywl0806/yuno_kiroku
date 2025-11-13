import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import {createNativeStackNavigator} from '@react-navigation/native-stack';
import {NavigationContainer} from '@react-navigation/native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import React from 'react';
import {ActivityIndicator, Modal, StyleSheet, View} from 'react-native';
import Toast from 'react-native-toast-message';

/* eslint-disable react/no-unstable-nested-components */
import Icon from 'react-native-vector-icons/FontAwesome';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';

import colors from './colors';
import {AuthProvider, useAuth} from './src/context/AuthContext';
import {MyTabBar} from './src/navigation/MyTabBar';
import {HomeScreen} from './src/screens/HomeScreen';
import {LoginScreen} from './src/screens/LoginScreen';
import {SettingScreen} from './src/screens/SettingScreen';
import {UploadScreen} from './src/screens/UploadScreen';

const Tab = createBottomTabNavigator();
const Stack = createNativeStackNavigator();

const styles = StyleSheet.create({
  header: {
    backgroundColor: colors.main_bg,
    height: 50,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
  },
  root: {
    flex: 1,
    backgroundColor: colors.main_bg,
  },
});
const queryClient = new QueryClient();

const MainTabs = () => {
  const [uploadModal, setUploadModal] = React.useState<boolean>(false);
  return (
    <>
      <View style={styles.header} />
      <View style={styles.root}>
        <Tab.Navigator
          tabBar={props => (
            <MyTabBar setUploadModal={setUploadModal} {...props} />
          )}
          screenOptions={{
            headerShown: false,
          }}
          initialRouteName="Upload">
          <Tab.Screen
            name="Home"
            component={HomeScreen}
            options={{
              tabBarIcon: ({color, size}) => (
                <MaterialCommunityIcons
                  name="home"
                  color={color}
                  size={size}
                />
              ),
            }}
          />

          <Tab.Screen
            name="Setting"
            component={SettingScreen}
            options={{
              tabBarIcon: ({color, size}) => (
                <Icon name="gear" color={color} size={size} />
              ),
            }}
          />
        </Tab.Navigator>
        <Modal
          animationType="slide"
          transparent={false}
          visible={uploadModal}
          onRequestClose={() => {
            setUploadModal(false);
          }}>
          <UploadScreen closeModal={() => setUploadModal(false)} />
        </Modal>
      </View>
    </>
  );
};

const AppNavigator = () => {
  const {isAuthenticated, isLoading} = useAuth();

  if (isLoading) {
    return (
      <View style={[styles.root, {justifyContent: 'center', alignItems: 'center'}]}>
        <ActivityIndicator size="large" color={colors.main_blue} />
      </View>
    );
  }

  return (
    <NavigationContainer>
      <Stack.Navigator screenOptions={{headerShown: false}}>
        {isAuthenticated ? (
          <Stack.Screen name="Main" component={MainTabs} />
        ) : (
          <Stack.Screen name="Login" component={LoginScreen} />
        )}
      </Stack.Navigator>
    </NavigationContainer>
  );
};

const App = () => {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <AppNavigator />
        </AuthProvider>
      </QueryClientProvider>
      <Toast />
    </>
  );
};
export default App;
