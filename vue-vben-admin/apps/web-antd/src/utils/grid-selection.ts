export function getGridSelectedIds<T extends { id?: string }>(grid: {
  getCheckboxRecords?: () => T[];
}) {
  const rows = grid?.getCheckboxRecords?.() ?? [];
  return rows
    .map((item) => item?.id)
    .filter((id): id is string => typeof id === 'string' && id.length > 0);
}
