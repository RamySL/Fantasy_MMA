export function getWinnerId(winnerValue) {
  if (winnerValue === null || winnerValue === undefined) {
    return null;
  }

  if (typeof winnerValue === "number") {
    return winnerValue;
  }

  // Le backend Go peut serialiser sql.NullInt64 sous la forme
  // { Int64: number, Valid: boolean }. On accepte aussi une future forme simple.
  if (typeof winnerValue === "object") {
    if (winnerValue.Valid === false || winnerValue.valid === false) {
      return null;
    }

    const value = winnerValue.Int64 ?? winnerValue.int64 ?? winnerValue.value;
    return value === undefined || value === null ? null : Number(value);
  }

  return Number(winnerValue) || null;
}

export function getFightWinnerName(fight) {
  const winnerId = getWinnerId(fight.winner_id);

  if (!winnerId) {
    return null;
  }

  if (fight.fighter1?.id === winnerId) {
    return fight.fighter1.full_name;
  }

  if (fight.fighter2?.id === winnerId) {
    return fight.fighter2.full_name;
  }

  return null;
}
