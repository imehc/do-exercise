import { useState } from 'react'
import { Button } from '~/components/ui/button'

export function Page3() {
  const [state, setState] = useState<number>(0)
  return (
    <div>
      <h1>Page3</h1>
      <div>{state}</div>
      <Button onClick={() => setState(prev => prev + 1)}>更新</Button>
    </div>
  )
}

export default Page3
