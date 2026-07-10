function pad(value: number) {
  return String(value).padStart(2, '0')
}

// Business dates must use the user's local calendar day. ISO strings are UTC
// based and would otherwise turn early-morning China time into the prior day.
export function formatLocalDate(value = new Date()) {
  return `${value.getFullYear()}-${pad(value.getMonth() + 1)}-${pad(value.getDate())}`
}
