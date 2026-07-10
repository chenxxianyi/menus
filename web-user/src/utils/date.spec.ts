import { describe, expect, it } from 'vitest'
import { formatLocalDate } from './date'

describe('formatLocalDate', () => {
  it('uses the local calendar date instead of the UTC ISO date', () => {
    const localMidnight = new Date(2026, 6, 10, 0, 30, 0)
    expect(formatLocalDate(localMidnight)).toBe('2026-07-10')
  })

  it('pads month and day values', () => {
    expect(formatLocalDate(new Date(2026, 0, 2, 12, 0, 0))).toBe('2026-01-02')
  })
})
