import {PhotoIdentifier} from '@react-native-camera-roll/camera-roll';
import {useQuery} from '@tanstack/react-query';
import React from 'react';

import {getLivePhotoMov} from '../../utils/getLivePhotoMov';
import {LivePhotoView} from './LivePhotoView';

type Props = {
  photo: PhotoIdentifier;
};

export const LivePhotoLocal = ({photo}: Props) => {
  const {data} = useQuery({
    queryKey: ['liveUrl', photo.node.id],
    queryFn: () => getLivePhotoMov(photo.node.id),
  });

  return <LivePhotoView photoUrl={photo.node.image.uri} liveUrl={data ?? ''} />;
};
