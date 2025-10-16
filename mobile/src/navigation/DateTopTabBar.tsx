import {MaterialTopTabBarProps} from '@react-navigation/material-top-tabs';
import React from 'react';
import {
  Animated,
  ScrollView,
  StyleSheet,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';

import colors from '../../colors';

type DateTopTabBarProps = MaterialTopTabBarProps;

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.main_bg,
    flex: 1,
    maxHeight: 50,
  },
  tabIndicator: {
    width: '100%',

    height: 1,
    backgroundColor: colors.main_blue,
  },
});

export const DateTopTabBar = ({
  state,
  descriptors,
  navigation,
  position,
}: DateTopTabBarProps) => {
  return (
    <ScrollView style={styles.container} horizontal className="flex border-2">
      {state.routes.map((route, index) => {
        const {options} = descriptors[route.key];
        const label =
          options.tabBarLabel !== undefined
            ? options.tabBarLabel
            : options.title !== undefined
            ? options.title
            : route.name;

        const isFocused = state.index === index;

        const onPress = () => {
          const event = navigation.emit({
            type: 'tabPress',
            target: route.key,
            canPreventDefault: true,
          });

          if (!isFocused && !event.defaultPrevented) {
            navigation.navigate(route.name, route.params);
          }
        };

        const onLongPress = () => {
          navigation.emit({
            type: 'tabLongPress',
            target: route.key,
          });
        };

        const inputRange = state.routes.map((_, i) => i);
        const opacity = position.interpolate({
          inputRange,
          outputRange: inputRange.map(i => (i === index ? 1 : 0.3)),
        });
        const opacity2 = position.interpolate({
          inputRange,
          outputRange: inputRange.map(i => (i === index ? 1 : 0)),
        });
        return (
          <TouchableWithoutFeedback
            key={route.key}
            accessibilityRole="button"
            accessibilityState={isFocused ? {selected: true} : {}}
            accessibilityLabel={options.tabBarAccessibilityLabel}
            testID={options.tabBarTestID}
            onPress={onPress}
            onLongPress={onLongPress}>
            <View className="flex w-7 justify-center items-center">
              <Animated.Text style={{opacity}}>
                {label.toString()}
              </Animated.Text>
              <Animated.View
                style={[
                  styles.tabIndicator,
                  {
                    opacity: opacity2,
                  },
                ]}
              />
            </View>
          </TouchableWithoutFeedback>
        );
      })}
    </ScrollView>
  );
};
