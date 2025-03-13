declare namespace globalThis {
  type UnwrapPromise<T> = T extends Promise<infer U> ? U : T;

  type RouteParams<T = unknown> = {
    params: Promise<T>
  }

  type RouteSearchParams<T = unknown> = {
    searchParams: Promise<T>
  }
}