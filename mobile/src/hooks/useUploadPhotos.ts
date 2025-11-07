import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';
import {useMutation} from '@tanstack/react-query';
import Toast from 'react-native-toast-message';

import {uploadLivePhoto} from '../services/uploadLivePhoto';
import {uploadPhoto} from '../services/uploadPhoto';
import {getLivePhotoMov} from '../utils/getLivePhotoMov';

export const useUploadPhotos = () => {
  const {mutateAsync, isPending} = useMutation({
    mutationKey: ['uploadPhotos'],
    mutationFn: async (photos: PhotoIdentifier[]) => {
      const process = [];

      try {
        for await (const photo of photos) {
          if (photo.node.subTypes.includes('PhotoLive')) {
            const liveUrl = await getLivePhotoMov(photo.node.id);

            process.push(uploadLivePhoto(photo, liveUrl));
            continue;
          }
          process.push(uploadPhoto(photo));
        }

        const results = await Promise.all(process);
        Toast.show({
          type: 'success',
          text1: 'アップロード成功！',
        });

        return results;
      } catch (error) {
        Toast.show({
          type: 'error',
          text1: 'アップロード失敗だよ〜',
          text2: 'ヨンウに連絡してね〜',
        });
        throw error;
      }
    },
  });

  return {uploadPhotos: mutateAsync, isPending};
};
