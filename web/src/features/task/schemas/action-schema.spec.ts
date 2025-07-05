import { isValidCron } from 'cron-validator'
import { describe, it, expect } from 'vitest'

const valid6 = '0 0 0 * * *'
const valid6b = '0,30 * * * * *'
const invalid = 'invalid cron'
const empty = ''

describe('cron-validator isValidCron (6 fields only)', () => {
  it('should pass with a valid 6-field cron expression (with seconds)', () => {
    expect(isValidCron(valid6, { seconds: true })).toBe(true)
  })

  it('should pass with a valid 6-field cron expression (every 30 seconds)', () => {
    expect(isValidCron(valid6b, { seconds: true })).toBe(true)
  })

  it('should fail with an invalid cron expression', () => {
    expect(isValidCron(invalid, { seconds: true })).toBe(false)
  })

  it('should fail with an empty cron expression', () => {
    expect(isValidCron(empty, { seconds: true })).toBe(false)
  })
})
