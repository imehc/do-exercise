import { z } from "zod";
import { paginationSchema } from "~/helper/schema";

export const postListSchema = paginationSchema.extend({
    name: z.string().optional()
})