declare namespace globalThis {
  type UnwrapPromise<T> = T extends Promise<infer U> ? U : T;
  
  type RouteParams = {
    params: Promise<unknown>
    searchParams: Promise<unknown>
  }
}