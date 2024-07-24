import { useState } from "react";
import {
  useMutation,
  useQueryClient,
  UseMutationOptions,
} from "@tanstack/react-query";
import axios, { AxiosError, Method } from "axios";

interface useRequestParams {
  url: string;
  method: Method;
  body?: any;
  onSuccess?: (value?: any) => void;
}

interface RequestError {
  message: string;
}

interface UseRequestReturnType {
  requestErrors: RequestError[] | null;
  sendRequest: (params: useRequestParams) => Promise<any>;
  isLoading: boolean;
}

interface AxiosErrorResponse {
  errors: RequestError[];
}

export default function useRequest(): UseRequestReturnType {
  const [requestErrors, setRequestErrors] = useState<RequestError[] | null>(
    null
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async ({ url, method, body }: useRequestParams) => {
      const response = await axios({
        url,
        method,
        data: body,
      });
      return response.data;
    },
    onSuccess: (data, variables) => {
      if (variables?.onSuccess) {
        variables.onSuccess(data);
      }
      queryClient.invalidateQueries();
    },
    onError: (error: AxiosError<AxiosErrorResponse>) => {
      if (error.response?.data.errors) {
        setRequestErrors(error.response.data.errors);
      } else {
        setRequestErrors([{ message: "Some Error Occurred" }]);
      }
    },
    onSettled: () => {
      setIsLoading(false);
    },
  } as UseMutationOptions<
    any,
    AxiosError<AxiosErrorResponse>,
    useRequestParams
  >);

  const sendRequest = async (params: useRequestParams): Promise<any> => {
    setRequestErrors(null);
    setIsLoading(true);
    const result = await mutation.mutateAsync(params);
    setIsLoading(false);
    return result;
  };

  return { requestErrors, sendRequest, isLoading };
}
