import type { Range, Salary } from './api'

export function formatDate(value: string | null): string {
  if (!value) return '—'

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleDateString('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

function years(count: number): string {
  return `${count} ${count === 1 ? 'year' : 'years'}`
}

export function formatExperience(range: Range): string {
  const { min, max } = range

  if (min !== null && max !== null) {
    return min === max ? years(min) : `${min}–${max} years`
  }
  if (min !== null) {
    return `${min}+ years`
  }
  if (max !== null) {
    return `Up to ${years(max)}`
  }
  return 'Not specified'
}

function money(value: number, currency: string | null): string {
  if (currency) {
    try {
      return new Intl.NumberFormat('en', {
        style: 'currency',
        currency,
        maximumFractionDigits: 0,
      }).format(value)
    } catch {
      return `${currency} ${new Intl.NumberFormat('en', { maximumFractionDigits: 0 }).format(value)}`
    }
  }

  return new Intl.NumberFormat('en', { maximumFractionDigits: 0 }).format(value)
}

export function formatSalary(salary: Salary): string {
  if (!salary.disclosed) {
    return 'Not disclosed'
  }

  const { min, max, currency, period } = salary
  let amount: string

  if (min !== null && max !== null && min !== max) {
    amount = `${money(min, currency)} – ${money(max, currency)}`
  } else {
    amount = money(min ?? max ?? 0, currency)
  }

  return period ? `${amount} / ${period}` : amount
}

export function titleCase(value: string): string {
  if (!value) return '—'

  return value
    .replaceAll('_', ' ')
    .split(' ')
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}
