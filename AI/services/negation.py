"""Shared, bounded negation handling for sentiment and urgency."""

NEGATION_WORDS = frozenset({"tidak", "bukan", "tak", "belum", "jangan", "tanpa", "tiada", "kurang"})
NEGATION_WINDOW = 3


def is_negated(tokens: list[str], index: int) -> bool:
    """Return whether a token is within the short scope of a negation."""
    start = max(0, index - NEGATION_WINDOW)
    return any(token in NEGATION_WORDS for token in tokens[start:index])


def effective_tokens(tokens: list[str]) -> list[str]:
    """Remove negated terms and negation markers for rule/model scoring."""
    return [
        token
        for index, token in enumerate(tokens)
        if token not in NEGATION_WORDS and not is_negated(tokens, index)
    ]
