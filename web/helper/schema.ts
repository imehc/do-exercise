import { z } from "zod";

/** 
 * 分页参数
 */
export const paginationSchema = z.object({
    page: z.coerce.number().catch(1).default(1),
    pageSize: z.coerce.number().catch(10).default(10),
})