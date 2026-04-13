import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useQuery } from '@tanstack/react-query'

export type KidFaceImg = {
  kid_id: number
  name: string
  face_img_url: string
  media_item_id: number
  taken_at_year: number
  taken_at_month: number
  birth_date: string
}

type KidFaceImgResponse = {
  kids: KidFaceImg[]
  fast_kids: KidFaceImg[]
}

/**
 * 월별 아이 얼굴 사진 조회
 * @param year 년도
 * @param month 월
 * @returns 아이 얼굴 사진 목록
 */
export const useGetKidFaceImgs = (year: number, month: number) => {
  return useQuery({
    queryKey: ['kid-face-imgs', year, month],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<KidFaceImgResponse>(API_ROUTES.KID.FACE_IMGS, {
        params: { year, month },
      })
      return response.data
    },
    staleTime: 1000 * 60 * 60 * 24,
  })
}
