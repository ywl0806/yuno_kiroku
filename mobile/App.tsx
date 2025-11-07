import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import {NavigationContainer} from '@react-navigation/native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import React from 'react';
import {Modal, StyleSheet, View} from 'react-native';
import Toast from 'react-native-toast-message';

/* eslint-disable react/no-unstable-nested-components */
import Icon from 'react-native-vector-icons/FontAwesome';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';

import colors from './colors';
import {MyTabBar} from './src/navigation/MyTabBar';
import {HomeScreen} from './src/screens/HomeScreen';
import {SettingScreen} from './src/screens/SettingScreen';
import {UploadScreen} from './src/screens/UploadScreen';

const Tab = createBottomTabNavigator();

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
const App = () => {
  const [uploadModal, setUploadModal] = React.useState<boolean>(false);
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <NavigationContainer>
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
        </NavigationContainer>
      </QueryClientProvider>
      <Toast />
    </>
  );
};
export default App;
