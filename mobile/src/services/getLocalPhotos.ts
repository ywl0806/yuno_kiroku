import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';

export type LocalPhotos = [string, PhotoIdentifier[]][];

export const groupPhotosByDate = (photos: PhotoIdentifier[]): LocalPhotos => {
  const grouped = photos.reduce((acc, photo) => {
    const date = new Date(photo.node.timestamp * 1000);
    const dateKey = date.toISOString().split('T')[0];
    const existing = acc.find(([key]) => key === dateKey);
    if (existing) {
      existing[1].push(photo);
    } else {
      acc.push([dateKey, [photo]]);
    }
    return acc;
  }, [] as LocalPhotos);
  return grouped;
};
