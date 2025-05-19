export const callDisabledTypes = new Map<true | false, string>([
  [
    true,
    'bg-destructive/10 dark:bg-destructive/50 text-destructive dark:text-primary border-destructive/10',
  ],
  [false, 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200'],
])

export const callMethodTypes = new Map<string, string>([
  [
    'GET',
    'bg-green-100/30 text-green-900 dark:text-green-200 border-green-200',
  ],
  [
    'POST',
    'bg-orange-100/30 text-orange-900 dark:text-orange-200 border-orange-200',
  ],
  ['PUT', 'bg-blue-100/30 text-blue-900 dark:text-blue-200 border-blue-200'],
  ['PATCH', 'bg-pink-100/30 text-pink-900 dark:text-pink-200 border-pink-200'],
  [
    'DELETE',
    'bg-destructive/10 dark:bg-destructive/50 text-destructive dark:text-primary border-destructive/10',
  ],
])
