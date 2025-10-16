import React from 'react';
import {Image, Text, View} from 'react-native';

import {livePhotoIcon} from '../../icons/livePhotoIcon';

type Props = {
  size?: 'small' | 'large';
};
export const LivePhotoBadge = ({size = 'small'}: Props) => {
  const badgeSize = size === 'small' ? 12 : 20;
  const textSize = size === 'small' ? 9 : 11;
  return (
    <View
      className={`absolute top-2 left-2 z-10 rounded-xl bg-white/70 py-[2px] ${
        size === 'small' ? 'px-[2px]' : 'px-[4px]'
      }  flex flex-row items-center justify-center border-[0.5px] border-light_gray`}
      style={{
        shadowColor: '#000',
        shadowOffset: {
          width: 1,
          height: 1,
        },

        shadowRadius: 2.22,
        elevation: 2,
      }}>
      <Image
        source={livePhotoIcon}
        style={{width: badgeSize, height: badgeSize}}
      />
      <Text className="text-main_black mx-[1px] " style={{fontSize: textSize}}>
        LIVE
      </Text>
    </View>
  );
};
