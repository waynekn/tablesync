import { z } from "zod";

export const spreadsheetInitSchema = z.object({
  title: z.string().min(1, "Title is required"),
  description: z.string().min(1, "Description is required"),
  deadline: z
    .string()
    .transform((val) => {
      if (val) {
        const padded = val.length === 16 ? val + ":00" : val; // adds seconds if missing
        return new Date(padded).toISOString(); // converts to full RFC3339 with 'Z'
      }
      return "";
    })
    .refine((val) => !isNaN(Date.parse(val)), {
      message: "Invalid date format",
    }),

  colTitles: z
    .array(z.string())
    .min(1, "At least one column title is required"),
});
