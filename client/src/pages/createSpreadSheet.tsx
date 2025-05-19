import { useState } from "react";
import { AxiosError } from "axios";
import { List } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { withAuthenticationRequired } from "react-oidc-context";

import { useApi } from "@/hooks/api";

import { Button } from "@/components/ui/button";
import ColTitlesList from "@/components/colTitlesList";

import { spreadsheetInitSchema } from "@/schema/spreadsheet";
import { SpreadsheetInit, SpreadsheetInitErr } from "@/types/spreadsheet";

function CreateSheet() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [columnInput, setColumnInput] = useState(""); // State for column input
  const api = useApi();
  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    setError,
    watch,
  } = useForm<SpreadsheetInit>({
    resolver: zodResolver(spreadsheetInitSchema),
    defaultValues: {
      title: "",
      description: "",
      deadline: "",
      colTitles: [],
    },
  });

  const { mutate, isPending } = useMutation({
    mutationFn: (data: SpreadsheetInit) => {
      return api.post<string>("/spreadsheet/create/", data);
    },
    onError: (error: AxiosError<SpreadsheetInitErr>) => {
      const err = error as AxiosError<SpreadsheetInitErr>;
      const errorData = err.response?.data;

      let hasFieldErrors = false;

      if (errorData) {
        if (errorData.title) {
          setError("title", { type: "server", message: errorData.title });
          hasFieldErrors = true;
        }
        if (errorData.description) {
          setError("description", {
            type: "server",
            message: errorData.description,
          });
          hasFieldErrors = true;
        }
        if (errorData.deadline) {
          setError("deadline", {
            type: "server",
            message: errorData.deadline,
          });
          hasFieldErrors = true;
        }
        if (errorData.colTitles) {
          setError("colTitles", {
            type: "server",
            message: errorData.colTitles,
          });
          hasFieldErrors = true;
        }
        if (errorData.error) {
          setError("root", {
            type: "server",
            message: errorData.error,
          });
          hasFieldErrors = true;
        }
      }
      // If no specific field errors, show fallback root error
      if (!hasFieldErrors) {
        let fallbackMessage = "An unknown error occurred.";
        if (err.message) {
          fallbackMessage = err.message;
        } else if (err.code === "ERR_NETWORK") {
          fallbackMessage = "Network error: Unable to reach the server.";
        }
        setError("root", {
          type: "server",
          message: fallbackMessage,
        });
      }
    },
  });

  const columns = watch("colTitles"); // Watch the columns field

  const handleAddColumn = () => {
    if (columnInput.trim()) {
      const updatedColumns = [...(columns || []), columnInput.trim()];
      setValue("colTitles", updatedColumns, { shouldValidate: true }); // Update form state
      setColumnInput("");
    }
  };

  const onSubmit: SubmitHandler<SpreadsheetInit> = (formData) => {
    if (isPending) return; // Prevent multiple submissions
    mutate(formData);
  };

  const rmHeader = (index: number) => {
    const updatedColumns = columns.filter((_, i) => i !== index);
    setValue("colTitles", updatedColumns, { shouldValidate: true });
  };

  return (
    <div className="h-screen w-screen flex flex-col md:flex-row bg-gray-100">
      <ColTitlesList
        columns={columns}
        toggleSidebar={() => setIsSidebarOpen((prev) => !prev)}
        isSidebarOpen={isSidebarOpen}
        rmHeader={rmHeader}
      />

      <main className=" flex-1 p-4 md:p-8 sm:text-sm md:text-base">
        <Button
          className="md:hidden mb-4"
          onClick={() => setIsSidebarOpen((prev) => !prev)}
        >
          <List /> View Column Titles
        </Button>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="w-full max-w-3xl mx-auto border-2 border-gray-300 rounded-lg bg-white shadow-lg p-8 space-y-6"
        >
          <div className="flex flex-col">
            <label htmlFor="title" className="font-medium mb-2">
              Sheet title
            </label>
            <input
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-zinc-900 focus:border-zinc-900 sm:text-sm"
              placeholder="Title"
              id="title"
              {...register("title")}
            />
            {errors.title && (
              <p className="text-red-500 text-sm">{errors.title.message}</p>
            )}
          </div>

          <div className="flex flex-col">
            <label htmlFor="description" className="font-medium mb-2">
              Description
            </label>
            <input
              type="text"
              placeholder="Description"
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-zinc-900 focus:border-zinc-900 sm:text-sm"
              id="description"
              {...register("description")}
            />
            {errors.description && (
              <p className="text-red-600 text-sm">
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="flex flex-col">
            <label htmlFor="deadline" className="font-medium mb-2">
              Deadline
            </label>
            <input
              placeholder="Deadline"
              type="datetime-local"
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-zinc-900 focus:border-zinc-900 sm:text-sm"
              id="deadline"
              {...register("deadline")}
            />
            {errors.deadline && (
              <p className="text-red-500 text-sm">{errors.deadline.message}</p>
            )}
          </div>

          <div className="flex flex-col">
            <label htmlFor="columnInput" className="font-medium mb-2">
              Column Titles
            </label>
            <div className="flex gap-2">
              <input
                placeholder="Add column title"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-zinc-900 focus:border-zinc-900 sm:text-sm"
                id="columnInput"
                value={columnInput}
                onChange={(e) => setColumnInput(e.target.value)}
              />
              <Button
                type="button"
                onClick={handleAddColumn}
                disabled={!columnInput.trim()}
              >
                Add
              </Button>
            </div>
            {errors.colTitles && (
              <p className="text-red-500 text-sm">{errors.colTitles.message}</p>
            )}
          </div>

          {errors.root && (
            <p className="text-red-500 text-sm">{errors.root.message}</p>
          )}
          <Button type="submit" variant="dark" disabled={isPending}>
            Create
          </Button>
        </form>
      </main>
    </div>
  );
}

export default withAuthenticationRequired(CreateSheet, {
  OnRedirecting: () => <div>Redirecting to the login page...</div>,
});
