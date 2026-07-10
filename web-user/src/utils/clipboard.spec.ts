import { afterEach, describe, expect, it, vi } from 'vitest'
import { copyText } from './clipboard'

describe('copyText', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uses the clipboard API in a secure context', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(window, 'isSecureContext', { value: true, configurable: true })
    Object.defineProperty(navigator, 'clipboard', { value: { writeText }, configurable: true })

    await expect(copyText('购物清单')).resolves.toBe(true)
    expect(writeText).toHaveBeenCalledWith('购物清单')
  })

  it('falls back to document copy when Clipboard API is unavailable', async () => {
    Object.defineProperty(window, 'isSecureContext', { value: false, configurable: true })
    Object.defineProperty(document, 'execCommand', { value: vi.fn(() => true), configurable: true })

    await expect(copyText('邀请码')).resolves.toBe(true)
  })
})
