import { describe, expect, it } from 'vitest'
import { makeTreeNodeDataState } from './tree-utils'

const tree = [
  {
    name: 'parent',
    value: 1,
    children: [{ name: 'child', value: 2 }],
  },
]

describe('makeTreeNodeDataState', () => {
  it('does not implicitly select children when only the parent is selected', () => {
    const [parent] = makeTreeNodeDataState(tree, [1])

    expect(parent.checked).toBe(true)
    expect(parent.children?.[0].checked).toBe(false)
  })

  it('marks a parent indeterminate when only a child is selected', () => {
    const [parent] = makeTreeNodeDataState(tree, [2])

    expect(parent.checked).toBe(false)
    expect(parent.hasSelectedChildren).toBe(true)
    expect(parent.children?.[0].checked).toBe(true)
  })
})
