import { MyAxios } from '../lib/myAxios'

type Range = {
  year: number
  month: number
}

export const getPhotosRange = async (): Promise<Range[]> => {
  const response = await MyAxios.get('/photo/range')

  return response.data
}
