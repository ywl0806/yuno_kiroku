import React, {useRef, useState} from 'react';
import {
  Animated,
  Image,
  SafeAreaView,
  StyleSheet,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import Video from 'react-native-video';

type Props = {
  photoUrl: string;
  liveUrl: string;
};

const styles = StyleSheet.create({
  container: {
    width: '100%',
    height: '100%',
  },
  image: {
    width: '100%',
    height: '100%',
  },
  live: {
    width: '100%',
    height: '100%',
  },
  animatedContainer: {
    width: '100%',
    height: '100%',
    position: 'absolute',

    top: 0,
    left: 0,
  },
});

export const LivePhotoView = ({photoUrl, liveUrl}: Props) => {
  const [viewLive, setViewLive] = useState(false);

  const photoAnim = useRef(new Animated.Value(1)).current;
  const liveAnim = useRef(new Animated.Value(0)).current;
  const videoRef = useRef<Video>(null);
  const photoView = () => {
    Animated.timing(photoAnim, {
      toValue: 1,
      duration: 500,
      useNativeDriver: true,
    }).start();
    Animated.timing(liveAnim, {
      toValue: 0,
      duration: 500,
      useNativeDriver: true,
    }).start();
  };

  const liveView = () => {
    Animated.timing(photoAnim, {
      toValue: 0,
      duration: 500,
      useNativeDriver: true,
    }).start();
    Animated.timing(liveAnim, {
      toValue: 1,
      duration: 500,
      useNativeDriver: true,
    }).start();
  };

  return (
    <TouchableWithoutFeedback
      style={styles.container}
      onLongPress={() => {
        videoRef.current?.seek(0);
        setViewLive(true);
        liveView();
      }}
      onPressOut={() => {
        setViewLive(false);
        photoView();
      }}>
      <SafeAreaView style={styles.container}>
        <Animated.View
          style={[
            styles.animatedContainer,
            {
              opacity: 1,
            },
          ]}>
          <View>
            <Image
              source={{uri: photoUrl}}
              resizeMode="contain"
              style={styles.image}
            />
          </View>
        </Animated.View>
        {liveUrl && (
          <Animated.View
            style={[
              styles.animatedContainer,
              {
                opacity: liveAnim,
              },
            ]}>
            <View>
              <Video
                ref={videoRef}
                source={{uri: liveUrl}}
                resizeMode="contain"
                paused={!viewLive}
                controls={false}
                style={styles.live}
              />
            </View>
          </Animated.View>
        )}
      </SafeAreaView>
    </TouchableWithoutFeedback>
  );
};
