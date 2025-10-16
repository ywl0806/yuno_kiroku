import {BottomTabBarProps} from '@react-navigation/bottom-tabs';
import React from 'react';
import {TouchableOpacity, TouchableWithoutFeedback, View} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialIcons';

import colors from '../../colors';

type MyTabBarProps = BottomTabBarProps & {
  setUploadModal: (value: boolean) => void;
};

export const MyTabBar = ({
  state,
  descriptors,
  navigation,

  setUploadModal,
}: MyTabBarProps) => {
  return (
    <View className="relative w-full ">
      <View className="absolute w-full h-0 left-0 z-10">
        <View
          className="fixed w-[60px] h-[60px] rounded-full mx-auto -mt-4 z-10 flex justify-center items-center"
          style={{backgroundColor: colors.main_bg}}>
          <TouchableOpacity
            onPress={() => {
              setUploadModal(true);
            }}
            className="w-[45px] h-[45px] bg-white rounded-full flex justify-center items-center">
            <Icon name="upload" color={colors.dark_gray} size={30} />
          </TouchableOpacity>
        </View>
      </View>
      <View className="flex flex-row h-[50px] justify-around mb-6 mt-2 px-5 bg-white mx-3 rounded-full items-center shadow-lg ">
        {state.routes.map((route, index) => {
          const {options} = descriptors[route.key];

          const isFocused = state.index === index;

          const onPress = () => {
            const event = navigation.emit({
              type: 'tabPress',
              target: route.key,
              canPreventDefault: true,
            });

            if (!isFocused && !event.defaultPrevented) {
              navigation.navigate(route.name);
            }
          };

          const onLongPress = () => {
            navigation.emit({
              type: 'tabLongPress',
              target: route.key,
            });
          };

          return (
            <TouchableWithoutFeedback
              key={route.key}
              accessibilityRole="button"
              accessibilityState={isFocused ? {selected: true} : {}}
              accessibilityLabel={options.tabBarAccessibilityLabel}
              testID={options.tabBarTestID}
              onPress={onPress}
              onLongPress={onLongPress}>
              <View className={`justify-center items-center rounded-full p-2 `}>
                {options.tabBarIcon &&
                  options.tabBarIcon({
                    focused: isFocused,
                    color: isFocused ? colors.main_blue : colors.dark_gray,
                    size: 22,
                  })}
                <View
                  className={`w-1 h-1  rounded-full mt-1 ${
                    isFocused ? 'bg-main_blue' : ''
                  }`}
                />
              </View>
            </TouchableWithoutFeedback>
          );
        })}
      </View>
    </View>
  );
};
