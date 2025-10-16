import {
  CameraRoll,
  PhotoIdentifier,
} from '@react-native-camera-roll/camera-roll';
import {useQuery} from '@tanstack/react-query';
import React, {useEffect, useMemo, useState} from 'react';
import {FlatList, StyleSheet, Text, TouchableOpacity, View} from 'react-native';
import Icon from 'react-native-vector-icons/FontAwesome';

import colors from '../../colors';
import {UploadImageDetailModal} from '../components/blocks/modals/UploadImageDetailModal';
import {SelectablePhoto} from '../components/elements/SeletablePhoto';
import {useUploadPhotos} from '../hooks/useUploadPhotos';
import {groupPhotosByDate} from '../services/getLocalPhotos';

const styles = StyleSheet.create({
  container: {
    width: '100%',
    height: '100%',

    paddingTop: 50,
  },
  item: {
    flex: 1,
    maxWidth: '33%',
    padding: 5,
    aspectRatio: 1,
  },
  image: {
    width: '100%',
    height: '100%',
  },
  header: {
    width: '100%',
    height: 50,
    top: 0,
  },
  modalStyle: {
    width: '100%',
    height: '100%',
    zIndex: 100,
  },
});
const ONE_PAGE = 20;

type Props = {
  closeModal: () => void;
};
export const UploadScreen = ({closeModal}: Props) => {
  const [photos, setPhotos] = React.useState<PhotoIdentifier[]>([]);
  const [first, setFirst] = useState<number>(ONE_PAGE);
  const [selectedPhotos, setSelectedPhotos] = useState<string[]>([]);

  const [detailView, setDetailView] = useState<boolean>(false);
  const [currentDetailPhotoIndex, setCurrentDetailPhotoIndex] =
    useState<number>(0);
  const currentDetailPhoto = useMemo(() => {
    if (photos.length > 0) {
      return photos[currentDetailPhotoIndex];
    }
    return null;
  }, [currentDetailPhotoIndex, photos]);

  const {data} = useQuery({
    queryKey: ['localPhotos', first],
    queryFn: () => CameraRoll.getPhotos({first}).then(r => r.edges),
  });
  const photoGroup = useMemo(() => {
    if (photos) {
      return groupPhotosByDate(photos);
    }
    return [];
  }, [photos]);
  useEffect(() => {
    if (data) {
      setPhotos(data);
    }
  }, [data]);

  const {uploadPhotos, isPending} = useUploadPhotos();
  const seletePhoto = (photo: PhotoIdentifier) => {
    if (selectedPhotos.includes(photo.node.id)) {
      setSelectedPhotos(prev => prev.filter(p => p !== photo.node.id));
    } else {
      setSelectedPhotos(prev => [...prev, photo.node.id]);
    }
  };
  const handleUploadPhotos = async (photosArg: PhotoIdentifier[]) => {
    uploadPhotos(photosArg).then(() => {
      setSelectedPhotos([]);
      closeModal();
    });
  };

  return (
    <View style={styles.container}>
      <View className="flex px-3">
        <TouchableOpacity onPress={closeModal}>
          <Icon name="angle-left" color={colors.dark_gray} size={40} />
        </TouchableOpacity>
      </View>
      <FlatList
        className="pb-28 px-1"
        onEndReached={() => {
          setFirst(first + ONE_PAGE);
        }}
        keyExtractor={item => item[0]}
        onEndReachedThreshold={0.01}
        data={photoGroup}
        renderItem={({item: photoGroupItem}) => (
          <View>
            <Text className="text-xl font-bold text-dark_gray mb-2 mt-4">
              {photoGroupItem[0]}
            </Text>
            <FlatList
              data={photoGroupItem[1]}
              renderItem={({item: photoItem}) => (
                <SelectablePhoto
                  onPress={() => {
                    setDetailView(true);
                    setCurrentDetailPhotoIndex(
                      photos.findIndex(p => p.node.id === photoItem.node.id),
                    );
                  }}
                  photo={photoItem}
                  isSeleted={selectedPhotos.includes(photoItem.node.id)}
                  seletePhoto={seletePhoto}
                />
              )}
              keyExtractor={item => item.node.image.uri}
              numColumns={3}
            />
          </View>
        )}
      />
      {/* image detail modal */}
      <View className="relative">
        <UploadImageDetailModal
          uploadPhotos={handleUploadPhotos}
          detailView={detailView}
          setDetailView={setDetailView}
          currentDetailPhoto={currentDetailPhoto}
          seletePhoto={seletePhoto}
          selectedPhotos={selectedPhotos}
          photos={photos}
          setCurrentDetailPhotoIndex={setCurrentDetailPhotoIndex}
          currentDetailPhotoIndex={currentDetailPhotoIndex}
          isLoading={isPending}
        />
      </View>
      <View className="w-full flex justify-center h-[60px] items-end px-2 mb-7">
        <View className="flex flex-row items-center justify-end w-full ">
          {isPending ? (
            <View>
              <Text className="font-bold text-dark_gray mr-2">
                アップロード中...
              </Text>
            </View>
          ) : (
            <TouchableOpacity
              className={`flex flex-row items-center justify-items-center rounded-lg p-3 z-10
            ${selectedPhotos.length > 0 ? 'bg-secondary_blue' : 'bg-light_gray'}
            `}
              onPress={() => {
                const p = photos.filter(photo =>
                  selectedPhotos.includes(photo.node.id),
                );
                handleUploadPhotos(p);
              }}
              disabled={selectedPhotos.length === 0}>
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
  );
};
