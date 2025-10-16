import {useQuery} from '@tanstack/react-query';
import {useMemo} from 'react';

import {getPhotosRange} from '../services/getPhotosRange';

export const usePhotosRange = () => {
  const {data, isFetched, refetch} = useQuery({
    queryKey: ['photosRange'],
    queryFn: () => {
      return getPhotosRange();
    },
  });

  const range = useMemo(() => {
    if (!data) return [];
    const sotedData = data.sort((a, b) => {
      if (a.year === b.year) {
        return b.month - a.month;
      }
      return b.year - a.year;
    });

    return sotedData;
  }, [data]);
  // const {data, isFetched} = useQuery({
  //   queryKey: ['firstPhoto'],
  //   queryFn: () => {
  //     return getFirstPhoto();
  //   },
  // });

  // const range = useMemo(() => {
  //   const rng = [];

  //   if (data) {
  //     const first = new Date(data.photo_created_at);
  //     const firstMonth = first.getMonth();
  //     const firstYear = first.getFullYear();

  //     const last = new Date();
  //     const lastMonth = last.getMonth();
  //     const lastYear = last.getFullYear();

  //     for (let year = lastYear; year >= firstYear; year--) {
  //       const startMonth = year === firstYear ? firstMonth : 0;
  //       const endMonth = year === lastYear ? lastMonth : 11;
  //       for (let month = endMonth; month >= startMonth; month--) {
  //         rng.push({year, month});
  //       }
  //     }
  //   }

  //   return rng;
  // }, [data]);

  return {range, isFetched, refetch};
};
