import { toast } from 'sonner'
import { ResponseError } from '~/do-exercise-api'

export function handleServerError(error: unknown) {
  // eslint-disable-next-line no-console
  console.log(error)

  let errMsg = 'Something went wrong!'

  if (
    error &&
    typeof error === 'object' &&
    'status' in error &&
    Number(error.status) === 204
  ) {
    errMsg = 'Content not found.'
  }

  if (error instanceof ResponseError) {
    errMsg = error.message
  }

  toast.error(errMsg)
}
