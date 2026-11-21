export interface ApiListOption {
  value: string;
  label: string;
}

/** Supports the supplied { lists: [{ data: '123' }] } response. */
export function mapListResponse(response: unknown): ApiListOption[] {
  const rows = Array.isArray(response)
    ? response
    : response && typeof response === 'object'
      ? 'lists' in response
        ? response.lists
        : 'data' in response
          ? response.data
          : null
      : null;
  if (!Array.isArray(rows)) throw new Error('Expected lists array');
  return rows.map((row: unknown) => {
    if (!row || typeof row !== 'object') throw new Error('Invalid list row');
    if (
      'value' in row &&
      'label' in row &&
      typeof row.value === 'string' &&
      typeof row.label === 'string'
    ) {
      return { value: row.value, label: row.label };
    }
    if ('data' in row && (typeof row.data === 'string' || typeof row.data === 'number')) {
      return { value: String(row.data), label: String(row.data) };
    }
    throw new Error('Expected data or value/label');
  });
}
