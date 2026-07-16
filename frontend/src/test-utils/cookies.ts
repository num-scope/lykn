import { removeCookie } from '@/lib/cookies'

/**
 * 移除 `document.cookie` 中可见的 Cookie，以隔离测试
 *
 * - 未传 `filter`：移除所有 Cookie
 * - `string`：只移除名称以该字符串开头的 Cookie
 * - `RegExp`：只移除 `filter.test(name)` 为 true 的 Cookie
 */
export function clearCookies(filter?: string | RegExp): void {
  if (typeof document === 'undefined') return

  for (const part of document.cookie.split(';')) {
    const name = part.split('=')[0]?.trim()
    if (!name) continue

    const shouldRemove =
      filter === undefined
        ? true
        : typeof filter === 'string'
          ? name.startsWith(filter)
          : filter.test(name)

    if (shouldRemove) {
      removeCookie(name)
    }
  }
}
