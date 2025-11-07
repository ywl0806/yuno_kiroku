import {useQuery} from '@tanstack/react-query';

import {getPhotosGroup} from '../services/getPhotosGroup';

export type UseGetPhotosProps = {
  year: number;
  month: number;
};
export const useGetPhotos = ({year, month}: UseGetPhotosProps) => {
  const {data, refetch, isFetched} = useQuery({
    queryKey: ['photos', year, month],
    queryFn: () => {
      const from = new Date(year, month - 1, 1);
      const to = new Date(year, month);
      return getPhotosGroup({from, to});
    },
    staleTime: 1000 * 60 * 30,
    // enabled: false,
  });

  return {photosGroups: data, refetch, isFetched};
};
