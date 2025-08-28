import { JSX } from 'react'
import { UseQueryResult } from '@tanstack/react-query'

/** 使用模式匹配处理 @link https://gabrielpichot.fr/blog/simplify-tanstack-react-query-state-handling-with-pattern-matching/ */
export function matchQueryStatus<T>(
  query: UseQueryResult<T>,
  options: {
    Loading: JSX.Element
    Errored: JSX.Element | ((error: unknown) => JSX.Element)
    Empty: JSX.Element
    Success: (
      data: UseQueryResult<T> & {
        data: NonNullable<UseQueryResult<T>['data']>
      }
    ) => JSX.Element
  }
): JSX.Element

export function matchQueryStatus<T>(
  query: UseQueryResult<T>,
  options: {
    Loading: JSX.Element
    Errored: JSX.Element | ((error: unknown) => JSX.Element)
    Success: (data: UseQueryResult<T>) => JSX.Element
  }
): JSX.Element

export function matchQueryStatus<T>(
  query: UseQueryResult<T>,
  {
    Loading,
    Errored,
    Empty,
    Success,
  }: {
    Loading: JSX.Element
    Errored: JSX.Element | ((error: unknown) => JSX.Element)
    Empty?: JSX.Element
    Success: (data: UseQueryResult<T>) => JSX.Element
  }
): JSX.Element {
  if (query.isLoading) {
    return Loading
  }

  if (query.isError) {
    if (typeof Errored === 'function') {
      return Errored(query.error)
    }
    return Errored
  }

  const isEmpty =
    query.data === undefined ||
    query.data === null ||
    (Array.isArray(query.data) && query.data.length === 0)

  if (isEmpty && Empty) {
    return Empty
  }

  return Success(query)
}
