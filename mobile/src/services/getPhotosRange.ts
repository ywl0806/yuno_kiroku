import {MyAxios} from '../lib/myAxios';

type Range = {
  year: number;
  month: number;
};

export const getPhotosRange = async (): Promise<Range[]> => {
  const url = `http://localhost:1323/api/photo/range`;
  const response = await MyAxios.get(url);
  return response.data;
};
