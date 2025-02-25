interface Pagination {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
  is_next: boolean;
  is_previous: boolean;
}

export interface ApiResponse<T> {
  data: T;
  metadata: {
    code: number;
    pagination: Pagination;
  };
  error: unknown;
}
