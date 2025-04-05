import { useState } from 'react'
import { Button } from '~/components/ui/button'

export function Page1() {
  const [state, setState] = useState<number>(0)
  return (
    <div className="overflow-ellipsis overflow-hidden whitespace-nowrap">
      Page1 Page1
      <div>{state}</div>
      <Button onClick={() => setState(prev => prev + 1)}>更新</Button>
    </div>
  )
}

export default Page1
