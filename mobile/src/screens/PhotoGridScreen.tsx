import {useFocusEffect} from '@react-navigation/native';
import React from 'react';

import {PhotoGrid} from '../components/blocks/PhotoGrid';
import {useTab} from '../context/tabContext';
import {useGetPhotos} from '../hooks/useGetPhotos';

type Props = {
  year: number;
  month: number;
};

export const PhotoGridScreen = ({
  route,
}: {
  route?: {
    params: Props;
  };
}) => {
  const {photosGroups, refetch} = useGetPhotos({
    year: route?.params?.year ?? new Date().getFullYear(),
    month: route?.params?.month ?? new Date().getMonth() + 1,
  });
  const {setYear, setMonth} = useTab();

  useFocusEffect(() => {
    // if (!isFetched) {
    //   refetch();
    // }

    if (route?.params?.year) {
      setYear(route.params.year);
    }
    if (route?.params?.month) {
      setMonth(route.params.month);
    }
  });

  return (
    <PhotoGrid photos={photosGroups?.[0].photos ?? []} refetch={refetch} />
  );
};
