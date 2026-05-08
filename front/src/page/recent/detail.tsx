import { BatchDetailContainer } from "@/feature/recent/components/batch-detail-container"
import { useParams } from "react-router-dom"

export const BatchDetailPage = () => {

    const { batchId } = useParams<{ batchId: string }>()

    return <BatchDetailContainer batchId={Number(batchId)} />
}