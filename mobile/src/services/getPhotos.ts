import {MyAxios} from '../lib/myAxios';
import {Photo} from '../types/photo';

export type GetPhotosOptions = {
  from: Date;
  to: Date;
};

export const getPhotos = async (
  options: GetPhotosOptions,
): Promise<Photo[]> => {
  const url = `http://localhost:1323/api/photo`;

  const response = await MyAxios.get(url, {
    params: {
      from: options.from.toISOString(),
      to: options.to.toISOString(),
    },
  });

  return response.data;
};
