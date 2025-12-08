import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';

import {MyAxios} from '../lib/myAxios';

export const uploadPhoto = async (photo: PhotoIdentifier, albumId: string) => {
  const formData = new FormData();

  formData.append('file', {
    uri: photo.node.image.uri,
    type: 'image/*',
    name: photo.node.image.filename,
  });

  const url = `http://localhost:1323/api/photo/upload?album_id=${albumId}`;

  const response = await MyAxios.post(url, formData);
  return response.data;
};
