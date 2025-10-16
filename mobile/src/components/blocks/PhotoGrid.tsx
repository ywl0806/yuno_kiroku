import MasonryList from '@react-native-seoul/masonry-list';
import React, {useMemo, useState} from 'react';
import {
  Image,
  StyleProp,
  StyleSheet,
  TouchableWithoutFeedback,
  View,
  ViewStyle,
} from 'react-native';

import colors from '../../../colors';
import {Photo} from '../../types/photo';
import {ViewDetailModal} from './modals/ViewDetailModal';

type Props = {
  photos: Photo[];
  refetch?: () => void;
};

const styles = StyleSheet.create({
  imageGrid: {
    display: 'flex',
    flex: 1,
    minHeight: 200,
  },
  image: {
    width: '100%',
    height: '100%',
  },

  imageBox: {
    alignSelf: 'stretch',
    padding: 1,
  },
});

const PhotoCard = ({
  photo,
  style,
  onPress,
}: {
  photo: Photo;
  onPress: () => void;
  style?: StyleProp<ViewStyle>;
}) => {
  const randomHeight = useMemo(() => Math.floor(Math.random() * 100) + 200, []);
  return (
    <TouchableWithoutFeedback onPress={onPress}>
      <View style={[style, styles.imageBox, {height: randomHeight}]}>
        <Image
          source={{uri: `http://localhost:1323/${photo.thumbnail_url}`}}
          style={styles.image}
        />
      </View>
    </TouchableWithoutFeedback>
  );
};
export const PhotoGrid = ({photos, refetch}: Props) => {
  const [detailView, setDetailView] = useState(false);
  const [currentDetailPhotoIndex, setCurrentDetailPhotoIndex] = useState(0);

  return (
    <View
      className="flex-1"
      style={{
        backgroundColor: colors.main_bg,
      }}>
      <MasonryList
        onRefresh={refetch}
        data={photos}
        renderItem={
          (({item, i}: {item: Photo; i: number}) => {
            return (
              <PhotoCard
                photo={item}
                onPress={() => {
                  setCurrentDetailPhotoIndex(i);
                  setDetailView(true);
                }}
              />
            );
          }) as any
        }
        numColumns={2}
      />
      <ViewDetailModal
        detailView={detailView}
        setDetailView={setDetailView}
        photos={photos}
        setCurrentDetailPhotoIndex={setCurrentDetailPhotoIndex}
        currentDetailPhotoIndex={currentDetailPhotoIndex}
      />
    </View>
  );
};
