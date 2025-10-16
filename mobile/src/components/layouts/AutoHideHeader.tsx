import React from 'react';
import {Animated} from 'react-native';

import colors from '../../../colors';
import {DropDownYearProps, DropdownYear} from '../blocks/DropdownYear';

type AnimatedAppBarProps = DropDownYearProps & {
  translateY: Animated.AnimatedInterpolation<string | number>;
};

export const AutoHideHeader = ({
  translateY,

  ...props
}: AnimatedAppBarProps) => {
  return (
    <Animated.View
      style={{
        flexDirection: 'row',
        justifyContent: 'space-around',
        backgroundColor: colors.main_bg,
        width: '100%',

        //for animation
        transform: [{translateY: translateY}],
        position: 'absolute',
        top: 0,
        right: 0,
        left: 0,
        elevation: 4,
        zIndex: 1,
      }}>
      <DropdownYear {...props} />
    </Animated.View>
  );
};
