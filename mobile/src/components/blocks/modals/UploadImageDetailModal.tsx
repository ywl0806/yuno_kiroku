import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';
import React from 'react';
import {
  Image,
  Modal,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import Swiper from 'react-native-swiper';
import Icon from 'react-native-vector-icons/FontAwesome';

import colors from '../../../../colors';
import {LivePhotoBadge} from '../../elements/LivePhotoBadge';
import {LivePhotoLocal} from '../LivePhotoLocal';

type Props = {
  detailView: boolean;
  setDetailView: (value: boolean) => void;
  currentDetailPhoto: PhotoIdentifier | null;
  seletePhoto: (photo: PhotoIdentifier) => void;
  selectedPhotos: string[];
  photos: PhotoIdentifier[];
  setCurrentDetailPhotoIndex: (value: number) => void;
  currentDetailPhotoIndex: number;
  uploadPhotos: (photos: PhotoIdentifier[]) => Promise<any> | any;
  isLoading: boolean;
};

const styles = StyleSheet.create({
  image: {
    width: '100%',
    height: '100%',
  },
});

export const UploadImageDetailModal = ({
  detailView,
  setDetailView,
  currentDetailPhoto,
  seletePhoto,
  selectedPhotos,
  photos,
  setCurrentDetailPhotoIndex,
  currentDetailPhotoIndex,
  uploadPhotos,
  isLoading,
}: Props) => {
  return (
    <Modal
      animationType="fade"
      transparent={true}
      visible={detailView}
      onRequestClose={() => {
        setDetailView(false);
      }}>
      <View className="h-full w-full pt-[55px] bg-main_bg">
        <View className="flex flex-row justify-between items-center mx-3">
          <TouchableOpacity
            className="rounded-full bg-light_gray/50 p-2"
            onPress={() => {
              setDetailView(false);
            }}>
            <Icon name="close" color={colors.dark_gray} size={20} />
          </TouchableOpacity>
          <TouchableWithoutFeedback
            className=""
            onPress={() => {
              if (currentDetailPhoto) seletePhoto(currentDetailPhoto);
            }}>
            <View
              className={`w-7 h-7 absolute rounded-full bg-white/70 z-10 right-0 flex justify-center items-center ${
                selectedPhotos.includes(currentDetailPhoto?.node.id ?? '')
                  ? 'border-main bg-white'
                  : 'border-slate-500/50'
              } border-[3px] transition-all duration-500`}>
              {selectedPhotos.includes(currentDetailPhoto?.node.id ?? '') && (
                <Icon name="check" color={colors.dark_gray} />
              )}
            </View>
          </TouchableWithoutFeedback>
        </View>
        <Swiper
          loop={false}
          horizontal
          showsPagination={false}
          onIndexChanged={index => {
            setCurrentDetailPhotoIndex(index);
          }}
          index={currentDetailPhotoIndex}>
          {photos.map((photo, index) => (
            <View key={index} className="p-1">
              <View
                style={{
                  position: 'absolute',
                }}>
                {photo.node.subTypes.includes('PhotoLive') && (
                  <LivePhotoBadge size="large" />
                )}
              </View>
              {photo.node.subTypes.includes('PhotoLive') ? (
                <LivePhotoLocal photo={photo} />
              ) : (
                <Image
                  source={{uri: photo.node.image.uri}}
                  resizeMode="contain"
                  style={styles.image}
                />
              )}
            </View>
          ))}
        </Swiper>
        <View className="w-full flex justify-center h-[60px] items-end px-2 mb-7">
          <View className="flex flex-row items-center justify-end w-full ">
            {isLoading ? (
              <View>
                <Text className="font-bold text-dark_gray mr-2">
                  アップロード中...
                </Text>
              </View>
            ) : (
              <TouchableOpacity
                className={`flex flex-row items-center justify-items-center rounded-lg p-3 z-10
            ${
              selectedPhotos.length > 0 || isLoading
                ? 'bg-secondary_blue'
                : 'bg-light_gray'
            }
            `}
                onPress={() => {
                  const p = photos.filter(photo =>
                    selectedPhotos.includes(photo.node.id),
                  );
                  uploadPhotos(p);
                }}
                disabled={selectedPhotos.length === 0 || isLoading}>
                <Text className="font-bold text-dark_gray mr-2">
                  {`${selectedPhotos.length} 個の画像をアップロードする`}
                </Text>
                <Icon
                  name="upload"
                  color={
                    selectedPhotos.length > 0
                      ? colors.main_blue
                      : colors.dark_gray
                  }
                  size={20}
                />
              </TouchableOpacity>
            )}
          </View>
        </View>
      </View>
    </Modal>
  );
};
