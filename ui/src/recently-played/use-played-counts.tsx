import { useQuery} from "@tanstack/react-query"
import { readFileFromData } from "./parse-local-data";

export function usePlayedCounts(userId: string) {
    return useQuery({
        queryKey: ['recently-played-counts'],
        queryFn: () => loadData(userId),
        staleTime: 30*60*1000,
    })
}

async function loadData(userId: string) {
   return await readFileFromData(userId); 
}


